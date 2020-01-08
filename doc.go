/*
package pipeline is a library intended to help build dynamic request pipelines

Why?

Let's imagine you'd like to wrap an existing Restful API which uses a special format of JSON messages to communicate.
The problem is that you have to transform the message on the way to the API as well as transform the response and the way back.
To solve the problem you would usually have to write custom code to wrap such API.

This library is intended to write a simple pipeline in JSON format to solve this problem without writing any code.
As the format to define a pipeline is JSON it's very easy to create/configure a pipeline from a JavaScript environment, e.g. a browser.

Example:
1. user makes a request with a JSON body
2. the first pipeline step transforms the JSON body
3. the second pipeline step makes a HTTP request with the transformed JSON body
4. the third pipeline step takes the JSON response from the request and transforms it again
5. the client gets back the transformed JSON, not needing to know that the intermediate formats of the JSON got transformed

*/
package pipeline
