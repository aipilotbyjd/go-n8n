# n8n Core Nodes - Developer Implementation Guide

## Table of Contents
1. [Node Architecture & Base Structure](#node-architecture--base-structure)
2. [Core Node Types Implementation](#core-node-types-implementation)
3. [Trigger Nodes](#trigger-nodes)
4. [Data Processing Nodes](#data-processing-nodes)
5. [Logic & Control Flow Nodes](#logic--control-flow-nodes)
6. [Integration Nodes](#integration-nodes)
7. [Utility Nodes](#utility-nodes)

---

## Node Architecture & Base Structure

### Base Node Interface
```typescript
interface INode {
  id: string;                    // Unique node instance ID
  name: string;                   // User-defined name
  type: string;                   // Node type identifier
  typeVersion: number;            // Node version
  position: [number, number];     // Canvas position [x, y]
  disabled?: boolean;             // Skip execution if true
  notes?: string;                 // User notes
  credentials?: INodeCredentials; // Attached credentials
  parameters: INodeParameters;    // Node configuration
  webhookId?: string;            // For webhook nodes
  externalHooks?: IExternalHooks;
}

interface INodeType {
  description: INodeTypeDescription;
  execute?(this: IExecuteFunctions): Promise<INodeExecutionData[][]>;
  trigger?(this: ITriggerFunctions): Promise<ITriggerResponse | undefined>;
  webhook?(this: IWebhookFunctions): Promise<IWebhookResponseData>;
  poll?(this: IPollFunctions): Promise<INodeExecutionData[][] | null>;
}

interface INodeTypeDescription {
  displayName: string;
  name: string;
  group: string[];               // ['input', 'output', 'transform']
  version: number;
  description: string;
  defaults: INodeTypeDefaults;
  inputs: string[];               // ['main'] for regular, [] for trigger
  outputs: string[];              // ['main'] for single, ['main', 'main'] for multiple
  credentials?: INodeCredentialDescription[];
  properties: INodeProperties[];  // Configuration fields
  webhooks?: IWebhookDescription[];
  polling?: boolean;
  trigger?: boolean;
  subtitle?: string;
}

interface INodeExecutionData {
  json: IDataObject;              // JSON data
  binary?: IBinaryKeyData;        // Binary attachments
  error?: Error;                  // Error if failed
  pairedItem?: IPairedItemData;   // Source tracking
}

interface IExecuteFunctions {
  getInputData(): INodeExecutionData[];
  getNodeParameter(parameterName: string, itemIndex: number): any;
  getWorkflow(): IWorkflow;
  getNode(): INode;
  getCredentials(type: string): Promise<ICredentials>;
  getExecutionId(): string;
  getTimezone(): string;
  getWorkflowStaticData(type: 'global' | 'node'): IDataObject;
  helpers: {
    request(options: IRequestOptions): Promise<any>;
    httpRequest(options: IHttpRequestOptions): Promise<any>;
    prepareBinaryData(buffer: Buffer, fileName?: string, mimeType?: string): Promise<IBinaryData>;
    getBinaryDataBuffer(propertyName: string, inputIndex?: number): Promise<Buffer>;
    copyBinaryFile(propertyName: string, destinationPath: string): Promise<string>;
  };
}
```

### Node Execution Flow
```typescript
class NodeExecutor {
  async execute(node: INode, inputData: INodeExecutionData[][]): Promise<INodeExecutionData[][]> {
    // 1. Validate inputs
    if (!this.validateInputs(node, inputData)) {
      throw new NodeOperationError('Invalid inputs');
    }

    // 2. Prepare execution context
    const context = this.createExecutionContext(node, inputData);
    
    // 3. Load credentials if needed
    if (node.credentials) {
      context.credentials = await this.loadCredentials(node.credentials);
    }
    
    // 4. Execute node
    try {
      const nodeType = this.nodeTypes.getByName(node.type);
      
      // Check execution type
      if (nodeType.trigger) {
        return await this.executeTrigger(nodeType, context);
      } else if (nodeType.poll) {
        return await this.executePoll(nodeType, context);
      } else if (nodeType.webhook) {
        return await this.executeWebhook(nodeType, context);
      } else {
        return await nodeType.execute.call(context);
      }
    } catch (error) {
      // 5. Handle errors
      if (node.continueOnFail) {
        return this.wrapError(error, inputData);
      }
      throw error;
    }
  }
}
```

---

## Core Node Types Implementation

### 1. Start Node
```typescript
class StartNode implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'Start',
    name: 'start',
    group: ['input'],
    version: 1,
    description: 'Workflow starting point',
    defaults: {
      name: 'Start',
      color: '#553399',
    },
    inputs: [],
    outputs: ['main'],
    properties: []
  };

  async execute(this: IExecuteFunctions): Promise<INodeExecutionData[][]> {
    // Start node passes through manual execution data
    const items = this.getInputData();
    
    if (items.length === 0) {
      // Return empty item if no input data
      items.push({ json: {} });
    }
    
    return this.prepareOutputData(items);
  }
}
```

### 2. Schedule Trigger Node
```typescript
class ScheduleTrigger implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'Schedule Trigger',
    name: 'scheduleTrigger',
    group: ['trigger'],
    version: 1,
    description: 'Triggers workflow on schedule',
    defaults: {
      name: 'Schedule Trigger',
      color: '#00FF00',
    },
    inputs: [],
    outputs: ['main'],
    trigger: true,
    properties: [
      {
        displayName: 'Trigger Type',
        name: 'triggerType',
        type: 'options',
        options: [
          { name: 'Cron', value: 'cron' },
          { name: 'Interval', value: 'interval' },
        ],
        default: 'interval',
      },
      {
        displayName: 'Cron Expression',
        name: 'cronExpression',
        type: 'string',
        displayOptions: {
          show: { triggerType: ['cron'] }
        },
        default: '0 * * * *',
        description: 'Cron expression for scheduling',
      },
      {
        displayName: 'Interval',
        name: 'interval',
        type: 'number',
        displayOptions: {
          show: { triggerType: ['interval'] }
        },
        default: 60,
        description: 'Interval in seconds',
      }
    ]
  };

  async trigger(this: ITriggerFunctions): Promise<ITriggerResponse | undefined> {
    const triggerType = this.getNodeParameter('triggerType') as string;
    
    if (triggerType === 'cron') {
      const cronExpression = this.getNodeParameter('cronExpression') as string;
      const cronJob = new CronJob(cronExpression, async () => {
        this.emit([this.helpers.returnJsonArray({ 
          timestamp: new Date().toISOString() 
        })]);
      });
      
      cronJob.start();
      
      // Return function to stop the cron job
      async function closeFunction() {
        cronJob.stop();
      }
      
      return {
        closeFunction,
      };
    } else {
      const interval = this.getNodeParameter('interval') as number;
      const intervalId = setInterval(() => {
        this.emit([this.helpers.returnJsonArray({ 
          timestamp: new Date().toISOString() 
        })]);
      }, interval * 1000);
      
      async function closeFunction() {
        clearInterval(intervalId);
      }
      
      return {
        closeFunction,
      };
    }
  }
}
```

### 3. Webhook Node
```typescript
class WebhookNode implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'Webhook',
    name: 'webhook',
    group: ['trigger'],
    version: 1,
    description: 'Creates HTTP endpoint',
    defaults: {
      name: 'Webhook',
      color: '#885577',
    },
    inputs: [],
    outputs: ['main'],
    webhooks: [
      {
        name: 'default',
        httpMethod: 'POST',
        responseMode: 'onReceived',
        path: '={{$parameter["path"]}}',
        restartWebhook: true,
      }
    ],
    properties: [
      {
        displayName: 'Path',
        name: 'path',
        type: 'string',
        default: 'webhook',
        required: true,
        description: 'Webhook endpoint path',
      },
      {
        displayName: 'HTTP Method',
        name: 'httpMethod',
        type: 'options',
        options: [
          { name: 'GET', value: 'GET' },
          { name: 'POST', value: 'POST' },
          { name: 'PUT', value: 'PUT' },
          { name: 'DELETE', value: 'DELETE' },
          { name: 'HEAD', value: 'HEAD' },
          { name: 'PATCH', value: 'PATCH' },
        ],
        default: 'POST',
      },
      {
        displayName: 'Response Mode',
        name: 'responseMode',
        type: 'options',
        options: [
          { name: 'On Received', value: 'onReceived' },
          { name: 'Last Node', value: 'lastNode' },
        ],
        default: 'onReceived',
      },
      {
        displayName: 'Response Code',
        name: 'responseCode',
        type: 'number',
        default: 200,
      },
      {
        displayName: 'Response Headers',
        name: 'responseHeaders',
        type: 'json',
        default: '{}',
      }
    ]
  };

  async webhook(this: IWebhookFunctions): Promise<IWebhookResponseData> {
    const req = this.getRequestObject();
    const resp = this.getResponseObject();
    const headers = this.getHeaderData();
    const params = this.getParamsData();
    const query = this.getQueryData();
    const body = this.getBodyData();
    
    const returnData: INodeExecutionData[] = [{
      json: {
        headers,
        params,
        query,
        body,
        method: req.method,
        url: req.url,
      }
    }];
    
    // Handle binary data (file uploads)
    if (req.files && Array.isArray(req.files)) {
      const binaryData: IBinaryKeyData = {};
      for (const file of req.files) {
        const binaryPropertyName = file.fieldname || 'data';
        binaryData[binaryPropertyName] = await this.helpers.prepareBinaryData(
          file.buffer,
          file.originalname,
          file.mimetype
        );
      }
      returnData[0].binary = binaryData;
    }
    
    const responseMode = this.getNodeParameter('responseMode') as string;
    
    if (responseMode === 'onReceived') {
      const responseCode = this.getNodeParameter('responseCode') as number;
      const responseHeaders = this.getNodeParameter('responseHeaders') as object;
      
      return {
        webhookResponse: {
          status: responseCode,
          headers: responseHeaders,
          body: { success: true }
        },
        workflowData: [returnData],
      };
    }
    
    // For 'lastNode' mode, continue workflow and respond at the end
    return {
      workflowData: [returnData],
    };
  }
}
```

### 4. HTTP Request Node
```typescript
class HttpRequestNode implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'HTTP Request',
    name: 'httpRequest',
    group: ['transform'],
    version: 3,
    description: 'Makes HTTP requests',
    defaults: {
      name: 'HTTP Request',
      color: '#0033AA',
    },
    inputs: ['main'],
    outputs: ['main'],
    credentials: [
      {
        name: 'httpBasicAuth',
        required: false,
      },
      {
        name: 'httpBearerTokenAuth',
        required: false,
      },
      {
        name: 'httpOAuth2Api',
        required: false,
      }
    ],
    properties: [
      {
        displayName: 'Method',
        name: 'method',
        type: 'options',
        options: [
          { name: 'GET', value: 'GET' },
          { name: 'POST', value: 'POST' },
          { name: 'PUT', value: 'PUT' },
          { name: 'PATCH', value: 'PATCH' },
          { name: 'DELETE', value: 'DELETE' },
          { name: 'HEAD', value: 'HEAD' },
          { name: 'OPTIONS', value: 'OPTIONS' },
        ],
        default: 'GET',
      },
      {
        displayName: 'URL',
        name: 'url',
        type: 'string',
        default: '',
        required: true,
      },
      {
        displayName: 'Authentication',
        name: 'authentication',
        type: 'options',
        options: [
          { name: 'None', value: 'none' },
          { name: 'Basic Auth', value: 'basicAuth' },
          { name: 'Bearer Token', value: 'bearerToken' },
          { name: 'OAuth2', value: 'oAuth2' },
        ],
        default: 'none',
      },
      {
        displayName: 'Send Query Parameters',
        name: 'sendQuery',
        type: 'boolean',
        default: false,
      },
      {
        displayName: 'Query Parameters',
        name: 'queryParameters',
        type: 'fixedCollection',
        displayOptions: {
          show: { sendQuery: [true] }
        },
        typeOptions: {
          multipleValues: true,
        },
        default: {},
        options: [
          {
            name: 'parameter',
            displayName: 'Parameter',
            values: [
              {
                displayName: 'Name',
                name: 'name',
                type: 'string',
                default: '',
              },
              {
                displayName: 'Value',
                name: 'value',
                type: 'string',
                default: '',
              }
            ]
          }
        ]
      },
      {
        displayName: 'Send Headers',
        name: 'sendHeaders',
        type: 'boolean',
        default: false,
      },
      {
        displayName: 'Headers',
        name: 'headers',
        type: 'fixedCollection',
        displayOptions: {
          show: { sendHeaders: [true] }
        },
        typeOptions: {
          multipleValues: true,
        },
        default: {},
        options: [
          {
            name: 'header',
            displayName: 'Header',
            values: [
              {
                displayName: 'Name',
                name: 'name',
                type: 'string',
                default: '',
              },
              {
                displayName: 'Value',
                name: 'value',
                type: 'string',
                default: '',
              }
            ]
          }
        ]
      },
      {
        displayName: 'Send Body',
        name: 'sendBody',
        type: 'boolean',
        default: false,
        displayOptions: {
          show: {
            method: ['POST', 'PUT', 'PATCH', 'DELETE'],
          },
        },
      },
      {
        displayName: 'Body Content Type',
        name: 'bodyContentType',
        type: 'options',
        displayOptions: {
          show: {
            sendBody: [true],
          },
        },
        options: [
          { name: 'JSON', value: 'json' },
          { name: 'Form Data', value: 'form' },
          { name: 'Form URL Encoded', value: 'urlencoded' },
          { name: 'Raw', value: 'raw' },
          { name: 'Binary', value: 'binary' },
        ],
        default: 'json',
      },
      {
        displayName: 'Body',
        name: 'body',
        type: 'json',
        displayOptions: {
          show: {
            sendBody: [true],
            bodyContentType: ['json'],
          },
        },
        default: '{}',
      },
      {
        displayName: 'Options',
        name: 'options',
        type: 'collection',
        placeholder: 'Add Option',
        default: {},
        options: [
          {
            displayName: 'Ignore Response Code',
            name: 'ignoreResponseCode',
            type: 'boolean',
            default: false,
            description: 'Succeeds also when status code is not 2xx',
          },
          {
            displayName: 'Redirect Policy',
            name: 'redirect',
            type: 'options',
            options: [
              { name: 'Follow', value: 'follow' },
              { name: 'Error', value: 'error' },
              { name: 'Manual', value: 'manual' },
            ],
            default: 'follow',
          },
          {
            displayName: 'Response Format',
            name: 'responseFormat',
            type: 'options',
            options: [
              { name: 'JSON', value: 'json' },
              { name: 'Text', value: 'text' },
              { name: 'Binary', value: 'binary' },
            ],
            default: 'json',
          },
          {
            displayName: 'Timeout',
            name: 'timeout',
            type: 'number',
            default: 60000,
            description: 'Timeout in milliseconds',
          },
          {
            displayName: 'Retry on Fail',
            name: 'retry',
            type: 'number',
            default: 0,
            description: 'Number of retries',
          },
          {
            displayName: 'Batch Size',
            name: 'batchSize',
            type: 'number',
            default: 1,
            description: 'Process items in batches',
          },
          {
            displayName: 'Batch Interval',
            name: 'batchInterval',
            type: 'number',
            default: 1000,
            description: 'Milliseconds between batches',
          },
          {
            displayName: 'Pagination',
            name: 'pagination',
            type: 'boolean',
            default: false,
          },
          {
            displayName: 'Pagination Mode',
            name: 'paginationMode',
            displayOptions: {
              show: { pagination: [true] }
            },
            type: 'options',
            options: [
              { name: 'Off', value: 'off' },
              { name: 'Update URL', value: 'url' },
              { name: 'Update Body', value: 'body' },
              { name: 'Response Contains Next URL', value: 'responseContainsNextUrl' },
            ],
            default: 'off',
          }
        ]
      }
    ]
  };

  async execute(this: IExecuteFunctions): Promise<INodeExecutionData[][]> {
    const items = this.getInputData();
    const returnData: INodeExecutionData[] = [];
    
    const method = this.getNodeParameter('method', 0) as string;
    const authentication = this.getNodeParameter('authentication', 0) as string;
    
    // Process each input item
    for (let i = 0; i < items.length; i++) {
      try {
        const url = this.getNodeParameter('url', i) as string;
        
        // Build request options
        const requestOptions: IRequestOptions = {
          method,
          url,
          json: true,
          gzip: true,
          rejectUnauthorized: !this.getNodeParameter('options.allowUnauthorizedCerts', i, false),
          timeout: this.getNodeParameter('options.timeout', i, 60000) as number,
        };
        
        // Add authentication
        if (authentication !== 'none') {
          await this.addAuthentication(requestOptions, authentication);
        }
        
        // Add query parameters
        if (this.getNodeParameter('sendQuery', i)) {
          const queryParameters = this.getNodeParameter('queryParameters.parameter', i, []) as IDataObject[];
          requestOptions.qs = {};
          for (const param of queryParameters) {
            requestOptions.qs[param.name as string] = param.value;
          }
        }
        
        // Add headers
        if (this.getNodeParameter('sendHeaders', i)) {
          const headers = this.getNodeParameter('headers.header', i, []) as IDataObject[];
          requestOptions.headers = {};
          for (const header of headers) {
            requestOptions.headers[header.name as string] = header.value;
          }
        }
        
        // Add body
        if (this.getNodeParameter('sendBody', i, false)) {
          const bodyContentType = this.getNodeParameter('bodyContentType', i) as string;
          
          switch (bodyContentType) {
            case 'json':
              requestOptions.body = this.getNodeParameter('body', i, {});
              break;
            case 'form':
              requestOptions.form = this.getNodeParameter('body', i, {});
              break;
            case 'urlencoded':
              requestOptions.form = this.getNodeParameter('body', i, {});
              requestOptions.headers = {
                ...requestOptions.headers,
                'Content-Type': 'application/x-www-form-urlencoded',
              };
              break;
            case 'raw':
              requestOptions.body = this.getNodeParameter('body', i, '');
              requestOptions.json = false;
              break;
            case 'binary':
              const binaryPropertyName = this.getNodeParameter('binaryPropertyName', i) as string;
              const binaryData = await this.helpers.getBinaryDataBuffer(binaryPropertyName, i);
              requestOptions.body = binaryData;
              requestOptions.json = false;
              break;
          }
        }
        
        // Handle pagination
        const pagination = this.getNodeParameter('options.pagination', i, false) as boolean;
        if (pagination) {
          returnData.push(...await this.handlePagination(requestOptions, i));
        } else {
          // Make request
          const response = await this.helpers.request(requestOptions);
          
          // Process response
          const responseFormat = this.getNodeParameter('options.responseFormat', i, 'json') as string;
          
          if (responseFormat === 'binary') {
            const binaryData = await this.helpers.prepareBinaryData(
              Buffer.from(response),
              'data',
              'application/octet-stream'
            );
            returnData.push({
              json: { success: true },
              binary: { data: binaryData },
            });
          } else {
            returnData.push({
              json: responseFormat === 'json' ? response : { data: response },
            });
          }
        }
      } catch (error) {
        const ignoreResponseCode = this.getNodeParameter('options.ignoreResponseCode', i, false) as boolean;
        
        if (ignoreResponseCode && error.statusCode) {
          returnData.push({
            json: {
              error: error.message,
              statusCode: error.statusCode,
              body: error.response?.body,
            },
          });
        } else {
          throw error;
        }
      }
    }
    
    return [returnData];
  }
  
  private async addAuthentication(options: IRequestOptions, type: string): Promise<void> {
    switch (type) {
      case 'basicAuth':
        const basicCredentials = await this.getCredentials('httpBasicAuth');
        options.auth = {
          user: basicCredentials.user as string,
          pass: basicCredentials.password as string,
        };
        break;
      case 'bearerToken':
        const tokenCredentials = await this.getCredentials('httpBearerTokenAuth');
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${tokenCredentials.token}`,
        };
        break;
      case 'oAuth2':
        const oAuth2Options = await this.getCredentials('httpOAuth2Api');
        // OAuth2 token refresh logic would go here
        break;
    }
  }
  
  private async handlePagination(
    baseOptions: IRequestOptions, 
    itemIndex: number
  ): Promise<INodeExecutionData[]> {
    const results: INodeExecutionData[] = [];
    const paginationMode = this.getNodeParameter('options.paginationMode', itemIndex) as string;
    const maxPages = this.getNodeParameter('options.maxPages', itemIndex, 100) as number;
    
    let currentPage = 0;
    let hasNextPage = true;
    let nextUrl = baseOptions.url;
    
    while (hasNextPage && currentPage < maxPages) {
      const requestOptions = { ...baseOptions, url: nextUrl };
      
      if (paginationMode === 'body') {
        requestOptions.body = {
          ...requestOptions.body,
          page: currentPage,
        };
      }
      
      const response = await this.helpers.request(requestOptions);
      results.push({ json: response });
      
      // Check for next page
      if (paginationMode === 'responseContainsNextUrl') {
        nextUrl = response.next || response.nextUrl || response.next_page_url;
        hasNextPage = !!nextUrl;
      } else if (paginationMode === 'url') {
        // Update page number in URL
        const url = new URL(nextUrl);
        url.searchParams.set('page', (currentPage + 1).toString());
        nextUrl = url.toString();
        hasNextPage = response.length > 0;
      } else {
        hasNextPage = response.length > 0;
      }
      
      currentPage++;
    }
    
    return results;
  }
}
```

### 5. Function Node
```typescript
class FunctionNode implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'Function',
    name: 'function',
    group: ['transform'],
    version: 1,
    description: 'Execute custom JavaScript code',
    defaults: {
      name: 'Function',
      color: '#FF6600',
    },
    inputs: ['main'],
    outputs: ['main'],
    properties: [
      {
        displayName: 'JavaScript Code',
        name: 'functionCode',
        type: 'string',
        typeOptions: {
          alwaysOpenEditWindow: true,
          editor: 'code',
          rows: 10,
        },
        default: `// Code here will run once per input item
// Access the current item with: $json, $binary
// Return data with: return { json: {}, binary: {} }

// Example:
const item = $json;
item.processed = true;
item.timestamp = new Date().toISOString();

return {
  json: item
};`,
        description: 'JavaScript code to execute',
      }
    ]
  };

  async execute(this: IExecuteFunctions): Promise<INodeExecutionData[][]> {
    const items = this.getInputData();
    const returnData: INodeExecutionData[] = [];
    const functionCode = this.getNodeParameter('functionCode', 0) as string;
    
    // Create sandbox context
    const sandbox = {
      // Built-in modules
      Buffer,
      console,
      Math,
      Date,
      JSON,
      Array,
      Object,
      String,
      Number,
      Boolean,
      RegExp,
      Error,
      Promise,
      
      // n8n specific helpers
      $input: {
        all: () => items,
        first: () => items[0],
        last: () => items[items.length - 1],
        item: (index: number) => items[index],
      },
      $items: (nodeName?: string, outputIndex?: number) => {
        if (!nodeName) return items;
        return this.getInputData(outputIndex || 0);
      },
      $node: this.getNode(),
      $workflow: this.getWorkflow(),
      $execution: {
        id: this.getExecutionId(),
        mode: this.getMode(),
        resumeUrl: this.getRestartUrl(),
      },
      $evaluateExpression: (expression: string, itemIndex: number) => {
        return this.helpers.evaluateExpression(expression, itemIndex);
      },
      $getWorkflowStaticData: (type: 'global' | 'node') => {
        return this.getWorkflowStaticData(type);
      },
      $now: DateTime.now(),
      $today: DateTime.now().startOf('day'),
      
      // Helper functions
      helpers: {
        // HTTP request helper
        request: async (options: IRequestOptions) => {
          return await this.helpers.request(options);
        },
        
        // Binary data helpers
        prepareBinaryData: async (buffer: Buffer, fileName?: string, mimeType?: string) => {
          return await this.helpers.prepareBinaryData(buffer, fileName, mimeType);
        },
        getBinaryDataBuffer: async (propertyName: string, inputIndex?: number) => {
          return await this.helpers.getBinaryDataBuffer(propertyName, inputIndex);
        },
        
        // Utility functions
        md5: (text: string) => crypto.createHash('md5').update(text).digest('hex'),
        sha256: (text: string) => crypto.createHash('sha256').update(text).digest('hex'),
        base64Encode: (text: string) => Buffer.from(text).toString('base64'),
        base64Decode: (text: string) => Buffer.from(text, 'base64').toString(),
        uuid: () => uuidv4(),
      }
    };
    
    // Process each item
    for (let itemIndex = 0; itemIndex < items.length; itemIndex++) {
      const item = items[itemIndex];
      
      // Item-specific context
      const itemContext = {
        ...sandbox,
        $json: item.json,
        $binary: item.binary || {},
        $itemIndex: itemIndex,
        $runIndex: this.getExecutionData().runIndex,
      };
      
      try {
        // Create function from code
        const fn = new Function(
          ...Object.keys(itemContext),
          `
          return (async function() {
            ${functionCode}
          })();
          `
        );
        
        // Execute function with context
        const result = await fn.call(this, ...Object.values(itemContext));
        
        // Process result
        if (result === null || result === undefined) {
          // Skip item if null/undefined returned
          continue;
        } else if (Array.isArray(result)) {
          // Multiple items returned
          for (const resultItem of result) {
            returnData.push(this.normalizeItem(resultItem));
          }
        } else {
          // Single item returned
          returnData.push(this.normalizeItem(result));
        }
      } catch (error) {
        if (this.continueOnFail()) {
          returnData.push({
            json: { error: error.message },
            error,
          });
        } else {
          throw error;
        }
      }
    }
    
    return [returnData];
  }
  
  private normalizeItem(item: any): INodeExecutionData {
    if (typeof item === 'object' && item !== null) {
      if ('json' in item || 'binary' in item) {
        // Already in correct format
        return {
          json: item.json || {},
          binary: item.binary,
        };
      } else {
        // Plain object, treat as json
        return {
          json: item,
        };
      }
    } else {
      // Primitive value
      return {
        json: { value: item },
      };
    }
  }
}
```

### 6. IF Node (Conditional)
```typescript
class IfNode implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'IF',
    name: 'if',
    group: ['transform'],
    version: 1,
    description: 'Conditional logic - route items to different branches',
    defaults: {
      name: 'IF',
      color: '#408061',
    },
    inputs: ['main'],
    outputs: ['main', 'main'],
    outputNames: ['true', 'false'],
    properties: [
      {
        displayName: 'Conditions',
        name: 'conditions',
        type: 'fixedCollection',
        typeOptions: {
          multipleValues: true,
        },
        default: {},
        options: [
          {
            name: 'boolean',
            displayName: 'Boolean',
            values: [
              {
                displayName: 'Value 1',
                name: 'value1',
                type: 'string',
                default: '',
              },
              {
                displayName: 'Operation',
                name: 'operation',
                type: 'options',
                options: [
                  { name: 'Equal', value: 'equal' },
                  { name: 'Not Equal', value: 'notEqual' },
                ],
                default: 'equal',
              },
              {
                displayName: 'Value 2',
                name: 'value2',
                type: 'options',
                options: [
                  { name: 'true', value: true },
                  { name: 'false', value: false },
                ],
                default: true,
              }
            ]
          },
          {
            name: 'number',
            displayName: 'Number',
            values: [
              {
                displayName: 'Value 1',
                name: 'value1',
                type: 'number',
                default: 0,
              },
              {
                displayName: 'Operation',
                name: 'operation',
                type: 'options',
                options: [
                  { name: 'Smaller', value: 'smaller' },
                  { name: 'Smaller Equal', value: 'smallerEqual' },
                  { name: 'Equal', value: 'equal' },
                  { name: 'Not Equal', value: 'notEqual' },
                  { name: 'Larger', value: 'larger' },
                  { name: 'Larger Equal', value: 'largerEqual' },
                ],
                default: 'equal',
              },
              {
                displayName: 'Value 2',
                name: 'value2',
                type: 'number',
                default: 0,
              }
            ]
          },
          {
            name: 'string',
            displayName: 'String',
            values: [
              {
                displayName: 'Value 1',
                name: 'value1',
                type: 'string',
                default: '',
              },
              {
                displayName: 'Operation',
                name: 'operation',
                type: 'options',
                options: [
                  { name: 'Contains', value: 'contains' },
                  { name: 'Not Contains', value: 'notContains' },
                  { name: 'Ends With', value: 'endsWith' },
                  { name: 'Equal', value: 'equal' },
                  { name: 'Not Equal', value: 'notEqual' },
                  { name: 'Is Empty', value: 'isEmpty' },
                  { name: 'Is Not Empty', value: 'isNotEmpty' },
                  { name: 'Regex', value: 'regex' },
                  { name: 'Starts With', value: 'startsWith' },
                ],
                default: 'equal',
              },
              {
                displayName: 'Value 2',
                name: 'value2',
                type: 'string',
                displayOptions: {
                  hide: {
                    operation: ['isEmpty', 'isNotEmpty'],
                  },
                },
                default: '',
              }
            ]
          },
          {
            name: 'dateTime',
            displayName: 'Date & Time',
            values: [
              {
                displayName: 'Value 1',
                name: 'value1',
                type: 'dateTime',
                default: '',
              },
              {
                displayName: 'Operation',
                name: 'operation',
                type: 'options',
                options: [
                  { name: 'Occurred After', value: 'after' },
                  { name: 'Occurred Before', value: 'before' },
                ],
                default: 'after',
              },
              {
                displayName: 'Value 2',
                name: 'value2',
                type: 'dateTime',
                default: '',
              }
            ]
          }
        ]
      },
      {
        displayName: 'Combine',
        name: 'combineOperation',
        type: 'options',
        options: [
          { name: 'ALL', value: 'all', description: 'All conditions must be true' },
          { name: 'ANY', value: 'any', description: 'At least one condition must be true' },
        ],
        default: 'all',
        description: 'How to combine multiple conditions',
      }
    ]
  };

  async execute(this: IExecuteFunctions): Promise<INodeExecutionData[][]> {
    const items = this.getInputData();
    const combineOperation = this.getNodeParameter('combineOperation', 0) as string;
    
    const trueItems: INodeExecutionData[] = [];
    const falseItems: INodeExecutionData[] = [];
    
    // Process each item
    for (let itemIndex = 0; itemIndex < items.length; itemIndex++) {
      const item = items[itemIndex];
      
      // Get all conditions
      const conditions = this.getNodeParameter('conditions', itemIndex, {}) as IDataObject;
      const conditionResults: boolean[] = [];
      
      // Evaluate boolean conditions
      if (conditions.boolean) {
        for (const condition of conditions.boolean as IDataObject[]) {
          const value1 = this.evaluateExpression(condition.value1 as string, itemIndex);
          const value2 = condition.value2;
          const operation = condition.operation as string;
          
          let result = false;
          switch (operation) {
            case 'equal':
              result = value1 == value2;
              break;
            case 'notEqual':
              result = value1 != value2;
              break;
          }
          conditionResults.push(result);
        }
      }
      
      // Evaluate number conditions
      if (conditions.number) {
        for (const condition of conditions.number as IDataObject[]) {
          const value1 = parseFloat(this.evaluateExpression(condition.value1 as string, itemIndex));
          const value2 = parseFloat(condition.value2 as string);
          const operation = condition.operation as string;
          
          let result = false;
          switch (operation) {
            case 'smaller':
              result = value1 < value2;
              break;
            case 'smallerEqual':
              result = value1 <= value2;
              break;
            case 'equal':
              result = value1 === value2;
              break;
            case 'notEqual':
              result = value1 !== value2;
              break;
            case 'larger':
              result = value1 > value2;
              break;
            case 'largerEqual':
              result = value1 >= value2;
              break;
          }
          conditionResults.push(result);
        }
      }
      
      // Evaluate string conditions
      if (conditions.string) {
        for (const condition of conditions.string as IDataObject[]) {
          const value1 = String(this.evaluateExpression(condition.value1 as string, itemIndex));
          const value2 = String(condition.value2);
          const operation = condition.operation as string;
          
          let result = false;
          switch (operation) {
            case 'contains':
              result = value1.includes(value2);
              break;
            case 'notContains':
              result = !value1.includes(value2);
              break;
            case 'endsWith':
              result = value1.endsWith(value2);
              break;
            case 'equal':
              result = value1 === value2;
              break;
            case 'notEqual':
              result = value1 !== value2;
              break;
            case 'isEmpty':
              result = value1.length === 0;
              break;
            case 'isNotEmpty':
              result = value1.length > 0;
              break;
            case 'regex':
              result = new RegExp(value2).test(value1);
              break;
            case 'startsWith':
              result = value1.startsWith(value2);
              break;
          }
          conditionResults.push(result);
        }
      }
      
      // Evaluate dateTime conditions
      if (conditions.dateTime) {
        for (const condition of conditions.dateTime as IDataObject[]) {
          const value1 = new Date(this.evaluateExpression(condition.value1 as string, itemIndex));
          const value2 = new Date(condition.value2 as string);
          const operation = condition.operation as string;
          
          let result = false;
          switch (operation) {
            case 'after':
              result = value1 > value2;
              break;
            case 'before':
              result = value1 < value2;
              break;
          }
          conditionResults.push(result);
        }
      }
      
      // Combine results
      let finalResult: boolean;
      if (conditionResults.length === 0) {
        finalResult = true; // No conditions = true
      } else if (combineOperation === 'all') {
        finalResult = conditionResults.every(result => result === true);
      } else {
        finalResult = conditionResults.some(result => result === true);
      }
      
      // Route item to appropriate output
      if (finalResult) {
        trueItems.push(item);
      } else {
        falseItems.push(item);
      }
    }
    
    return [trueItems, falseItems];
  }
  
  private evaluateExpression(expression: string, itemIndex: number): any {
    return this.helpers.evaluateExpression(expression, itemIndex);
  }
}
```

### 7. Set Node (Data Manipulation)
```typescript
class SetNode implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'Set',
    name: 'set',
    group: ['input'],
    version: 2,
    description: 'Set fields on items',
    defaults: {
      name: 'Set',
      color: '#0000FF',
    },
    inputs: ['main'],
    outputs: ['main'],
    properties: [
      {
        displayName: 'Mode',
        name: 'mode',
        type: 'options',
        options: [
          { 
            name: 'Manually Map', 
            value: 'manual',
            description: 'Set specific fields via UI' 
          },
          { 
            name: 'Expression (JSON)', 
            value: 'expression',
            description: 'Write JavaScript/JSON expression' 
          },
        ],
        default: 'manual',
      },
      {
        displayName: 'Keep Only Set',
        name: 'keepOnlySet',
        type: 'boolean',
        default: false,
        description: 'Keep only fields that are set. Remove all other fields.',
      },
      {
        displayName: 'Fields to Set',
        name: 'fields',
        type: 'fixedCollection',
        typeOptions: {
          multipleValues: true,
        },
        displayOptions: {
          show: {
            mode: ['manual'],
          },
        },
        default: {},
        options: [
          {
            name: 'values',
            displayName: 'Values',
            values: [
              {
                displayName: 'Field Name',
                name: 'name',
                type: 'string',
                default: '',
                placeholder: 'e.g., email',
              },
              {
                displayName: 'Field Type',
                name: 'type',
                type: 'options',
                options: [
                  { name: 'String', value: 'string' },
                  { name: 'Number', value: 'number' },
                  { name: 'Boolean', value: 'boolean' },
                  { name: 'Date', value: 'dateTime' },
                  { name: 'Object', value: 'object' },
                  { name: 'Array', value: 'array' },
                ],
                default: 'string',
              },
              {
                displayName: 'Value',
                name: 'value',
                type: 'string',
                default: '',
              }
            ]
          }
        ]
      },
      {
        displayName: 'JSON Expression',
        name: 'jsonExpression',
        type: 'json',
        displayOptions: {
          show: {
            mode: ['expression'],
          },
        },
        default: '{\n  "field": "value"\n}',
      }
    ]
  };

  async execute(this: IExecuteFunctions): Promise<INodeExecutionData[][]> {
    const items = this.getInputData();
    const mode = this.getNodeParameter('mode', 0) as string;
    const keepOnlySet = this.getNodeParameter('keepOnlySet', 0) as boolean;
    const returnData: INodeExecutionData[] = [];
    
    for (let itemIndex = 0; itemIndex < items.length; itemIndex++) {
      const item = items[itemIndex];
      let newItem: IDataObject = {};
      
      if (!keepOnlySet) {
        // Keep existing fields
        newItem = { ...item.json };
      }
      
      if (mode === 'manual') {
        // Manual mode - set specific fields
        const fields = this.getNodeParameter('fields.values', itemIndex, []) as IDataObject[];
        
        for (const field of fields) {
          const fieldName = field.name as string;
          const fieldType = field.type as string;
          let fieldValue = field.value;
          
          // Parse value based on type
          switch (fieldType) {
            case 'number':
              fieldValue = parseFloat(fieldValue as string);
              break;
            case 'boolean':
              fieldValue = fieldValue === 'true' || fieldValue === true;
              break;
            case 'dateTime':
              fieldValue = new Date(fieldValue as string).toISOString();
              break;
            case 'object':
              try {
                fieldValue = JSON.parse(fieldValue as string);
              } catch (e) {
                fieldValue = {};
              }
              break;
            case 'array':
              try {
                fieldValue = JSON.parse(fieldValue as string);
                if (!Array.isArray(fieldValue)) {
                  fieldValue = [fieldValue];
                }
              } catch (e) {
                fieldValue = [];
              }
              break;
            default:
              fieldValue = String(fieldValue);
          }
          
          // Support nested field names (e.g., "user.email")
          this.setNestedValue(newItem, fieldName, fieldValue);
        }
      } else {
        // Expression mode - evaluate JSON expression
        const jsonExpression = this.getNodeParameter('jsonExpression', itemIndex) as string;
        
        try {
          const evaluatedExpression = this.helpers.evaluateExpression(jsonExpression, itemIndex);
          
          if (typeof evaluatedExpression === 'object' && evaluatedExpression !== null) {
            if (keepOnlySet) {
              newItem = evaluatedExpression as IDataObject;
            } else {
              newItem = { ...newItem, ...evaluatedExpression as IDataObject };
            }
          }
        } catch (error) {
          throw new Error(`Invalid JSON expression: ${error.message}`);
        }
      }
      
      returnData.push({
        json: newItem,
        binary: item.binary,
      });
    }
    
    return [returnData];
  }
  
  private setNestedValue(obj: IDataObject, path: string, value: any): void {
    const keys = path.split('.');
    let current = obj;
    
    for (let i = 0; i < keys.length - 1; i++) {
      const key = keys[i];
      if (!(key in current) || typeof current[key] !== 'object') {
        current[key] = {};
      }
      current = current[key] as IDataObject;
    }
    
    current[keys[keys.length - 1]] = value;
  }
}
```

### 8. Merge Node
```typescript
class MergeNode implements INodeType {
  description: INodeTypeDescription = {
    displayName: 'Merge',
    name: 'merge',
    group: ['transform'],
    version: 2,
    description: 'Merge multiple data streams',
    defaults: {
      name: 'Merge',
      color: '#00BBCC',
    },
    inputs: ['main', 'main'],
    outputs: ['main'],
    inputNames: ['Input 1', 'Input 2'],
    properties: [
      {
        displayName: 'Mode',
        name: 'mode',
        type: 'options',
        options: [
          {
            name: 'Append',
            value: 'append',
            description: 'Combine all items from all inputs',
          },
          {
            name: 'Merge By Index',
            value: 'mergeByIndex',
            description: 'Merge items with same index',
          },
          {
            name: 'Merge By Key',
            value: 'mergeByKey',
            description: 'Merge items by matching key field',
          },
          {
            name: 'Multiplex',
            value: 'multiplex',
            description: 'Create all possible combinations',
          },
          {
            name: 'Pass Through',
            value: 'passThrough',
            description: 'Output specific input unchanged',
          },
          {
            name: 'Wait',
            value: 'wait',
            description: 'Wait for all inputs to arrive',
          },
        ],
        default: 'append',
      },
      {
        displayName: 'Output',
        name: 'output',
        type: 'options',
        displayOptions: {
          show: {
            mode: ['passThrough'],
          },
        },
        options: [
          { name: 'Input 1', value: 'input1' },
          { name: 'Input 2', value: 'input2' },
        ],
        default: 'input1',
      },
      {
        displayName: 'Property Input 1',
        name: 'propertyName1',
        type: 'string',
        default: '',
        displayOptions: {
          show: {
            mode: ['mergeByKey'],
          },
        },
        description: 'Key field name in Input 1',
      },
      {
        displayName: 'Property Input 2',
        name: 'propertyName2',
        type: 'string',
        default: '',
        displayOptions: {
          show: {
            mode: ['mergeByKey'],
          },
        },
        description: 'Key field name in Input 2',
      },
      {
        displayName: 'Join',
        name: 'join',
        type: 'options',
        displayOptions: {
          show: {
            mode: ['mergeByKey'],
          },
        },
        options: [
          { name: 'Inner Join', value: 'inner' },
          { name: 'Left Join', value: 'left' },
          { name: 'Right Join', value: 'right' },
          { name: 'Outer Join', value: 'outer' },
        ],
        default: 'inner',
      },
      {
        displayName: 'Overwrite',
        name: 'overwrite',
        type: 'options',
        displayOptions: {
          show: {
            mode: ['mergeByIndex', 'mergeByKey'],
          },
        },
        options: [
          { name: 'Always', value: 'always' },
          { name: 'If Blank', value: 'blank' },
          { name: 'If Key Matches', value: 'key' },
        ],
        default: 'always',
      },
    ]
  };

  async execute(this: IExecuteFunctions): Promise<INodeExecutionData[][]> {
    const mode = this.getNodeParameter('mode', 0) as string;
    const returnData: INodeExecutionData[] = [];
    
    // Get data from both inputs
    const input1 = this.getInputData(0);
    const input2 = this.getInputData(1);
    
    switch (mode) {
      case 'append':
        // Simply combine all items
        returnData.push(...input1, ...input2);
        break;
        
      case 'mergeByIndex':
        // Merge items with same index
        const overwrite = this.getNodeParameter('overwrite', 0) as string;
        const maxLength = Math.max(input1.length, input2.length);
        
        for (let i = 0; i < maxLength; i++) {
          const item1 = input1[i];
          const item2 = input2[i];
          
          if (!item1) {
            returnData.push(item2);
          } else if (!item2) {
            returnData.push(item1);
          } else {
            const mergedItem: INodeExecutionData = {
              json: {},
              binary: {},
            };
            
            // Merge JSON data
            if (overwrite === 'always') {
              mergedItem.json = { ...item1.json, ...item2.json };
            } else if (overwrite === 'blank') {
              mergedItem.json = { ...item1.json };
              for (const key in item2.json) {
                if (!(key in item1.json) || item1.json[key] === '' || item1.json[key] === null) {
                  mergedItem.json[key] = item2.json[key];
                }
              }
            }
            
            // Merge binary data
            mergedItem.binary = {
              ...item1.binary,
              ...item2.binary,
            };
            
            returnData.push(mergedItem);
          }
        }
        break;
        
      case 'mergeByKey':
        // SQL-like join based on key field
        const propertyName1 = this.getNodeParameter('propertyName1', 0) as string;
        const propertyName2 = this.getNodeParameter('propertyName2', 0) as string;
        const join = this.getNodeParameter('join', 0) as string;
        
        // Build lookup map from input2
        const input2Map = new Map<any, INodeExecutionData[]>();
        for (const item of input2) {
          const key = this.getFieldValue(item.json, propertyName2);
          if (!input2Map.has(key)) {
            input2Map.set(key, []);
          }
          input2Map.get(key)!.push(item);
        }
        
        // Track matched items for outer joins
        const matchedInput2 = new Set<INodeExecutionData>();
        
        // Process input1 items
        for (const item1 of input1) {
          const key1 = this.getFieldValue(item1.json, propertyName1);
          const matches = input2Map.get(key1) || [];
          
          if (matches.length > 0) {
            // Found matches
            for (const item2 of matches) {
              matchedInput2.add(item2);
              returnData.push({
                json: { ...item1.json, ...item2.json },
                binary: { ...item1.binary, ...item2.binary },
              });
            }
          } else if (join === 'left' || join === 'outer') {
            // No match but include item1 (left join)
            returnData.push(item1);
          }
        }
        
        // Handle right and outer joins
        if (join === 'right' || join === 'outer') {
          for (const item2 of input2) {
            if (!matchedInput2.has(item2)) {
              returnData.push(item2);
            }
          }
        }
        break;
        
      case 'multiplex':
        // Create Cartesian product
        for (const item1 of input1) {
          for (const item2 of input2) {
            returnData.push({
              json: { ...item1.json, ...item2.json },
              binary: { ...item1.binary, ...item2.binary },
            });
          }
        }
        break;
        
      case 'passThrough':
        // Pass through specific input unchanged
        const output = this.getNodeParameter('output', 0) as string;
        if (output === 'input1') {
          returnData.push(...input1);
        } else {
          returnData.push(...input2);
        }
        break;
        
      case 'wait':
        // Wait for both inputs to have data
        if (input1.length > 0 && input2.length > 0) {
          returnData.push(...input1, ...input2);
        }
        break;
    }
    
    return [returnData];
  }
  
  private getFieldValue(obj: IDataObject, path: string): any {
    const keys = path.split('.');
    let value: any = obj;
    
    for (const key of keys) {
      if (value && typeof value === 'object' && key in value) {
        value = value[key];
      } else {
        return undefined;
      }
    }
    
    return value;
  }
}
```

## Implementation Patterns

### Error Handling Pattern
```typescript
class ErrorHandlingPattern {
  async executeWithRetry(
    fn: () => Promise<any>,
    options: {
      maxRetries: number;
      retryDelay: number;
      backoffMultiplier: number;
      onError?: (error: Error, attempt: number) => void;
    }
  ): Promise<any> {
    let lastError: Error;
    let delay = options.retryDelay;
    
    for (let attempt = 1; attempt <= options.maxRetries; attempt++) {
      try {
        return await fn();
      } catch (error) {
        lastError = error;
        
        if (options.onError) {
          options.onError(error, attempt);
        }
        
        if (attempt < options.maxRetries) {
          await new Promise(resolve => setTimeout(resolve, delay));
          delay *= options.backoffMultiplier;
        }
      }
    }
    
    throw lastError!;
  }
}
```

### Pagination Pattern
```typescript
class PaginationPattern {
  async *paginate(
    fetchFn: (page: number, limit: number) => Promise<any[]>,
    options: {
      pageSize: number;
      maxPages?: number;
    }
  ): AsyncGenerator<any[], void, unknown> {
    let page = 0;
    let hasMore = true;
    
    while (hasMore && (!options.maxPages || page < options.maxPages)) {
      const items = await fetchFn(page, options.pageSize);
      
      if (items.length > 0) {
        yield items;
        page++;
        hasMore = items.length === options.pageSize;
      } else {
        hasMore = false;
      }
    }
  }
}
```

### Batching Pattern
```typescript
class BatchingPattern {
  async processBatches<T, R>(
    items: T[],
    processFn: (batch: T[]) => Promise<R[]>,
    options: {
      batchSize: number;
      concurrency: number;
      delayBetweenBatches?: number;
    }
  ): Promise<R[]> {
    const results: R[] = [];
    const batches: T[][] = [];
    
    // Create batches
    for (let i = 0; i < items.length; i += options.batchSize) {
      batches.push(items.slice(i, i + options.batchSize));
    }
    
    // Process batches with concurrency control
    const queue = [...batches];
    const processing: Promise<void>[] = [];
    
    while (queue.length > 0 || processing.length > 0) {
      while (processing.length < options.concurrency && queue.length > 0) {
        const batch = queue.shift()!;
        
        const promise = processFn(batch).then(batchResults => {
          results.push(...batchResults);
          
          if (options.delayBetweenBatches) {
            return new Promise(resolve => 
              setTimeout(resolve, options.delayBetweenBatches)
            );
          }
        });
        
        processing.push(promise);
      }
      
      if (processing.length > 0) {
        await Promise.race(processing);
        processing.splice(
          processing.findIndex(p => p === await Promise.race(processing)),
          1
        );
      }
    }
    
    return results;
  }
}
```

## Node Development Best Practices

1. **Always validate inputs** - Check for required parameters and valid data types
2. **Handle errors gracefully** - Provide clear error messages and support continue-on-fail
3. **Support expressions** - Allow users to use n8n expressions in all string fields
4. **Preserve data lineage** - Include pairedItem references for data tracking
5. **Optimize for large datasets** - Use streaming and pagination where possible
6. **Document thoroughly** - Include descriptions, examples, and placeholder text
7. **Test edge cases** - Empty inputs, null values, large datasets, network failures
8. **Follow n8n conventions** - Consistent naming, colors, and UI patterns
9. **Support binary data** - Handle file uploads/downloads when relevant
10. **Implement proper cleanup** - Close connections, clear timeouts, free resources

---

This guide provides the complete technical foundation for implementing n8n core nodes in your clone. Each node example includes the full interface, execution logic, and common patterns used throughout n8n's architecture.
