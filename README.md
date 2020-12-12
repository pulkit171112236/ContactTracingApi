# ContactTracingApi
contact tracing api which accepts get request and post request in json

1) for adding a user
provide a http post request to url: /users

json body:
{
  "name": "<username>",
   "dob": "<dob>",
  "phone_num": "<phone number>",
  "email": "<email>"
}
  
//timestamp will be added automatically at host using time.Now()
