# Uisvc.DefaultApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**cancelRunIdPost**](DefaultApi.md#cancelRunIdPost) | **POST** /cancel/{run_id} | Cancel by Run ID
[**errorsGet**](DefaultApi.md#errorsGet) | **GET** /errors | Retrieve errors
[**logAttachIdGet**](DefaultApi.md#logAttachIdGet) | **GET** /log/attach/{id} | Attach to a running log
[**loggedinGet**](DefaultApi.md#loggedinGet) | **GET** /loggedin | Check logged in state
[**loginGet**](DefaultApi.md#loginGet) | **GET** /login | Log into the system
[**logoutGet**](DefaultApi.md#logoutGet) | **GET** /logout | Log out of the system
[**repositoriesCiAddOwnerRepoGet**](DefaultApi.md#repositoriesCiAddOwnerRepoGet) | **GET** /repositories/ci/add/{owner}/{repo} | Add a specific repository to CI.
[**repositoriesCiDelOwnerRepoGet**](DefaultApi.md#repositoriesCiDelOwnerRepoGet) | **GET** /repositories/ci/del/{owner}/{repo} | Removes a specific repository from CI.
[**repositoriesMyGet**](DefaultApi.md#repositoriesMyGet) | **GET** /repositories/my | Fetch all the writable repositories for the user.
[**repositoriesSubAddOwnerRepoGet**](DefaultApi.md#repositoriesSubAddOwnerRepoGet) | **GET** /repositories/sub/add/{owner}/{repo} | Subscribe to a repository running CI
[**repositoriesSubDelOwnerRepoGet**](DefaultApi.md#repositoriesSubDelOwnerRepoGet) | **GET** /repositories/sub/del/{owner}/{repo} | Unsubscribe from a repository
[**repositoriesSubscribedGet**](DefaultApi.md#repositoriesSubscribedGet) | **GET** /repositories/subscribed | List all subscribed repositories
[**repositoriesVisibleGet**](DefaultApi.md#repositoriesVisibleGet) | **GET** /repositories/visible | Fetch all the repositories the user can view.
[**runRunIdGet**](DefaultApi.md#runRunIdGet) | **GET** /run/{run_id} | Get a run by ID
[**runsCountGet**](DefaultApi.md#runsCountGet) | **GET** /runs/count | Count the runs
[**runsGet**](DefaultApi.md#runsGet) | **GET** /runs | Obtain the run list for the user
[**submitGet**](DefaultApi.md#submitGet) | **GET** /submit | Perform a manual submission to tinyCI
[**tasksCountGet**](DefaultApi.md#tasksCountGet) | **GET** /tasks/count | Count the Tasks
[**tasksGet**](DefaultApi.md#tasksGet) | **GET** /tasks | Obtain the task list optionally filtering by repository and sha.
[**tasksRunsIdCountGet**](DefaultApi.md#tasksRunsIdCountGet) | **GET** /tasks/runs/{id}/count | Count the runs corresponding to the task ID.
[**tasksRunsIdGet**](DefaultApi.md#tasksRunsIdGet) | **GET** /tasks/runs/{id} | Obtain the run list based on the task ID.
[**tasksSubscribedGet**](DefaultApi.md#tasksSubscribedGet) | **GET** /tasks/subscribed | Obtain the list of tasks that belong to repositories you are subscribed to.
[**tokenDelete**](DefaultApi.md#tokenDelete) | **DELETE** /token | Remove and reset your tinyCI access token
[**tokenGet**](DefaultApi.md#tokenGet) | **GET** /token | Get a tinyCI access token
[**userPropertiesGet**](DefaultApi.md#userPropertiesGet) | **GET** /user/properties | Get information about the current user


<a name="cancelRunIdPost"></a>
# **cancelRunIdPost**
> cancelRunIdPost(runId)

Cancel by Run ID

Cancel the run by ID; this will actually trickle back and cancel the whole task, since it can no longer succeed in any way. Please keep in mind to stop runs, runners must implement a cancel poller. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var runId = 56; // Number | The ID of the run to retrieve


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.cancelRunIdPost(runId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **runId** | **Number**| The ID of the run to retrieve | 

### Return type

null (empty response body)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="errorsGet"></a>
# **errorsGet**
> [UserError] errorsGet()

Retrieve errors

Server retrieves any errors the last call(s) have set for you.

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.errorsGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[UserError]**](UserError.md)

### Authorization

[session](../README.md#session)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="logAttachIdGet"></a>
# **logAttachIdGet**
> logAttachIdGet(id)

Attach to a running log

For a given ID, find the log and if it is running, attach to it and start receiving the latest content from it. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var id = 56; // Number | The ID of the run to retrieve the log for.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.logAttachIdGet(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Number**| The ID of the run to retrieve the log for. | 

### Return type

null (empty response body)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="loggedinGet"></a>
# **loggedinGet**
> &#39;String&#39; loggedinGet()

Check logged in state

Validate the logged-in status of the user. Validates the session cookie against the internal database. If the user is logged in, a JSON string of \&quot;true\&quot; will be sent; otherwise an oauth redirect url will be passed for calling out to by the client. 

### Example
```javascript
var Uisvc = require('uisvc');

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.loggedinGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

**&#39;String&#39;**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="loginGet"></a>
# **loginGet**
> loginGet(code, state)

Log into the system

Handle the server side of the oauth challenge. It is important to preserve the cookie jar after this call is made, as session cookies are used to manage many of the calls in this API. 

### Example
```javascript
var Uisvc = require('uisvc');

var apiInstance = new Uisvc.DefaultApi();

var code = "code_example"; // String | The code github sent back to us with the callback, we use it in the OAuth2 exchange to validate the request. 

var state = "state_example"; // String | The state (randomized string) we sent with the original link; this is echoed back to us so we can further identify the user. 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.loginGet(code, state, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **code** | **String**| The code github sent back to us with the callback, we use it in the OAuth2 exchange to validate the request.  | 
 **state** | **String**| The state (randomized string) we sent with the original link; this is echoed back to us so we can further identify the user.  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="logoutGet"></a>
# **logoutGet**
> logoutGet()

Log out of the system

Conveniently clears session cookies. You will need to login again. Does not clear oauth tokens. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.logoutGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

[session](../README.md#session)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="repositoriesCiAddOwnerRepoGet"></a>
# **repositoriesCiAddOwnerRepoGet**
> repositoriesCiAddOwnerRepoGet(owner, repo)

Add a specific repository to CI.

Generates a hook secret and populates the user&#39;s repository with it and the hook URL. Returns 200 on success, 500 + error message on failure, or if the repository has already been added to CI. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var owner = "owner_example"; // String | owner of the repository, first part of github repository name such as 'erikh' in 'erikh/foo' 

var repo = "repo_example"; // String | name of the repository, the second half of the github repository name such as 'foo' in 'erikh/foo'. 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.repositoriesCiAddOwnerRepoGet(owner, repo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **owner** | **String**| owner of the repository, first part of github repository name such as &#39;erikh&#39; in &#39;erikh/foo&#39;  | 
 **repo** | **String**| name of the repository, the second half of the github repository name such as &#39;foo&#39; in &#39;erikh/foo&#39;.  | 

### Return type

null (empty response body)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="repositoriesCiDelOwnerRepoGet"></a>
# **repositoriesCiDelOwnerRepoGet**
> repositoriesCiDelOwnerRepoGet(owner, repo)

Removes a specific repository from CI.

Will fail if not added to CI already; does not currently clear the hook. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var owner = "owner_example"; // String | owner of the repository, first part of github repository name such as 'erikh' in 'erikh/foo' 

var repo = "repo_example"; // String | name of the repository, the second half of the github repository name such as 'foo' in 'erikh/foo'. 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.repositoriesCiDelOwnerRepoGet(owner, repo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **owner** | **String**| owner of the repository, first part of github repository name such as &#39;erikh&#39; in &#39;erikh/foo&#39;  | 
 **repo** | **String**| name of the repository, the second half of the github repository name such as &#39;foo&#39; in &#39;erikh/foo&#39;.  | 

### Return type

null (empty response body)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="repositoriesMyGet"></a>
# **repositoriesMyGet**
> RepositoryList repositoriesMyGet()

Fetch all the writable repositories for the user.

Returns a types.RepositoryList for all the repos a user has write access to.

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.repositoriesMyGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**RepositoryList**](RepositoryList.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="repositoriesSubAddOwnerRepoGet"></a>
# **repositoriesSubAddOwnerRepoGet**
> repositoriesSubAddOwnerRepoGet(owner, repo)

Subscribe to a repository running CI

Subscribing makes that repo&#39;s queue items appear in your home view. Returns 200 on success, 500 + error on failure. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var owner = "owner_example"; // String | owner of the repository, first part of github repository name such as 'erikh' in 'erikh/foo' 

var repo = "repo_example"; // String | owner of the repository, first part of github repository name such as 'erikh' in 'erikh/foo' 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.repositoriesSubAddOwnerRepoGet(owner, repo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **owner** | **String**| owner of the repository, first part of github repository name such as &#39;erikh&#39; in &#39;erikh/foo&#39;  | 
 **repo** | **String**| owner of the repository, first part of github repository name such as &#39;erikh&#39; in &#39;erikh/foo&#39;  | 

### Return type

null (empty response body)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="repositoriesSubDelOwnerRepoGet"></a>
# **repositoriesSubDelOwnerRepoGet**
> repositoriesSubDelOwnerRepoGet(owner, repo)

Unsubscribe from a repository

Unsubscribing removes any existing subscription. Either way, if nothing broke, it returns 200. Otherwise it returns 500 and the error. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var owner = "owner_example"; // String | owner of the repository, first part of github repository name such as 'erikh' in 'erikh/foo' 

var repo = "repo_example"; // String | owner of the repository, first part of github repository name such as 'erikh' in 'erikh/foo' 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.repositoriesSubDelOwnerRepoGet(owner, repo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **owner** | **String**| owner of the repository, first part of github repository name such as &#39;erikh&#39; in &#39;erikh/foo&#39;  | 
 **repo** | **String**| owner of the repository, first part of github repository name such as &#39;erikh&#39; in &#39;erikh/foo&#39;  | 

### Return type

null (empty response body)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="repositoriesSubscribedGet"></a>
# **repositoriesSubscribedGet**
> RepositoryList repositoriesSubscribedGet()

List all subscribed repositories

Returns a types.RepositoryList of all the repos the user is subscribed to.

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.repositoriesSubscribedGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**RepositoryList**](RepositoryList.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="repositoriesVisibleGet"></a>
# **repositoriesVisibleGet**
> RepositoryList repositoriesVisibleGet()

Fetch all the repositories the user can view.

Returns a types.RepositoryList for all the repos a user has view access to.

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.repositoriesVisibleGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**RepositoryList**](RepositoryList.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="runRunIdGet"></a>
# **runRunIdGet**
> Run runRunIdGet(runId)

Get a run by ID

Retrieve a Run by ID; this will return the full Run object including all relationships.

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var runId = 56; // Number | The ID of the run to retrieve


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.runRunIdGet(runId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **runId** | **Number**| The ID of the run to retrieve | 

### Return type

[**Run**](Run.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="runsCountGet"></a>
# **runsCountGet**
> &#39;Number&#39; runsCountGet(opts)

Count the runs

Count the runs, optionally filtering by repository or repository+SHA. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var opts = { 
  'repository': "repository_example", // String | 
  'sha': "sha_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.runsCountGet(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **repository** | **String**|  | [optional] 
 **sha** | **String**|  | [optional] 

### Return type

**&#39;Number&#39;**

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="runsGet"></a>
# **runsGet**
> RunList runsGet(opts)

Obtain the run list for the user

List all the runs, optionally filtering by repository or repository+SHA. Pagination controls are available. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var opts = { 
  'page': 0, // Number | pagination control: what page to retrieve in the query.
  'perPage': 100, // Number | pagination control: how many items counts as a page.
  'repository': "repository_example", // String | optional; the repository name to get the tasks for.
  'sha': "sha_example" // String | optional; the sha to get the tasks for.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.runsGet(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **Number**| pagination control: what page to retrieve in the query. | [optional] [default to 0]
 **perPage** | **Number**| pagination control: how many items counts as a page. | [optional] [default to 100]
 **repository** | **String**| optional; the repository name to get the tasks for. | [optional] 
 **sha** | **String**| optional; the sha to get the tasks for. | [optional] 

### Return type

[**RunList**](RunList.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="submitGet"></a>
# **submitGet**
> submitGet(repository, sha, opts)

Perform a manual submission to tinyCI

This allows a user to push a job instead of pushing to git or filing a pull request to trigger a job. It is available on the tinyCI UI and CLI client. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var repository = "repository_example"; // String | the repository owner/repo to be tested.

var sha = "sha_example"; // String | the sha or branch to be tested

var opts = { 
  'all': true // Boolean | Run all tests instead of relying on diff selection to pick them.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.submitGet(repository, sha, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **repository** | **String**| the repository owner/repo to be tested. | 
 **sha** | **String**| the sha or branch to be tested | 
 **all** | **Boolean**| Run all tests instead of relying on diff selection to pick them. | [optional] 

### Return type

null (empty response body)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="tasksCountGet"></a>
# **tasksCountGet**
> &#39;Number&#39; tasksCountGet(opts)

Count the Tasks

Perform a full count of tasks that meet the filter criteria (which can be no filter) and return it as integer. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var opts = { 
  'repository': "repository_example", // String | optional; repository for filtering
  'sha': "sha_example" // String | optional; sha for filtering
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.tasksCountGet(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **repository** | **String**| optional; repository for filtering | [optional] 
 **sha** | **String**| optional; sha for filtering | [optional] 

### Return type

**&#39;Number&#39;**

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="tasksGet"></a>
# **tasksGet**
> TaskList tasksGet(opts)

Obtain the task list optionally filtering by repository and sha.

The tasks list returns a list of Task objects that correspond to the query. Each query may contain pagination or filtering rules to limit its contents. It is strongly recommended to look at the \&quot;count\&quot; equivalents for these endpoints so that you can implement pagination more simply. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var opts = { 
  'page': 0, // Number | pagination control: what page to retrieve in the query.
  'perPage': 100, // Number | pagination control: how many items counts as a page.
  'repository': "repository_example", // String | optional; the repository name to get the tasks for.
  'sha': "sha_example" // String | optional; the sha to get the tasks for.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.tasksGet(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **Number**| pagination control: what page to retrieve in the query. | [optional] [default to 0]
 **perPage** | **Number**| pagination control: how many items counts as a page. | [optional] [default to 100]
 **repository** | **String**| optional; the repository name to get the tasks for. | [optional] 
 **sha** | **String**| optional; the sha to get the tasks for. | [optional] 

### Return type

[**TaskList**](TaskList.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="tasksRunsIdCountGet"></a>
# **tasksRunsIdCountGet**
> &#39;Number&#39; tasksRunsIdCountGet(id)

Count the runs corresponding to the task ID.

Get the count of runs that correspond to the task ID. Returns an integer. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var id = 56; // Number | the ID of the Task.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.tasksRunsIdCountGet(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Number**| the ID of the Task. | 

### Return type

**&#39;Number&#39;**

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="tasksRunsIdGet"></a>
# **tasksRunsIdGet**
> RunList tasksRunsIdGet(id, opts)

Obtain the run list based on the task ID.

The queue list only contains: * stuff * other junk 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var id = 56; // Number | the ID of the Task

var opts = { 
  'page': 0, // Number | pagination control: what page to retrieve in the query.
  'perPage': 100 // Number | pagination control: how many items counts as a page.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.tasksRunsIdGet(id, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Number**| the ID of the Task | 
 **page** | **Number**| pagination control: what page to retrieve in the query. | [optional] [default to 0]
 **perPage** | **Number**| pagination control: how many items counts as a page. | [optional] [default to 100]

### Return type

[**RunList**](RunList.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="tasksSubscribedGet"></a>
# **tasksSubscribedGet**
> TaskList tasksSubscribedGet(opts)

Obtain the list of tasks that belong to repositories you are subscribed to.

This call implements basic pagination over the entire task corpus that intersects with your subscription list. It returns a list of tasks. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var opts = { 
  'page': 0, // Number | pagination control: what page to retrieve in the query.
  'perPage': 100 // Number | pagination control: how many items counts as a page.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.tasksSubscribedGet(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **Number**| pagination control: what page to retrieve in the query. | [optional] [default to 0]
 **perPage** | **Number**| pagination control: how many items counts as a page. | [optional] [default to 100]

### Return type

[**TaskList**](TaskList.md)

### Authorization

[session](../README.md#session), [token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="tokenDelete"></a>
# **tokenDelete**
> tokenDelete()

Remove and reset your tinyCI access token

The next GET /token will create a new one. This will just remove it. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: token
var token = defaultClient.authentications['token'];
token.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//token.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.tokenDelete(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="tokenGet"></a>
# **tokenGet**
> &#39;String&#39; tokenGet()

Get a tinyCI access token

This will allow you unfettered access to the system as your user that you request the token with. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.tokenGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

**&#39;String&#39;**

### Authorization

[session](../README.md#session)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="userPropertiesGet"></a>
# **userPropertiesGet**
> userPropertiesGet()

Get information about the current user

Get information about the current user, such as the username. 

### Example
```javascript
var Uisvc = require('uisvc');
var defaultClient = Uisvc.ApiClient.instance;

// Configure API key authorization: session
var session = defaultClient.authentications['session'];
session.apiKey = 'YOUR API KEY';
// Uncomment the following line to set a prefix for the API key, e.g. "Token" (defaults to null)
//session.apiKeyPrefix = 'Token';

var apiInstance = new Uisvc.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.userPropertiesGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

[session](../README.md#session)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

