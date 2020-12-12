# ContactTracingApi
#contact tracing api which accepts get request and post request in json

1) for creating a new user
provide a http post request to url: /users

json body:
{
  "name": "<username>",
   "dob": "<dob>",
  "phone_num": "<phone number>",
  "email": "<email>"
}  
#timestamp will be added automatically at host using time.Now()
#output will be _id by which the user is inserted in json format


2) for viewing all users
make a get request to url: /users
# output will be listing of all users with their whole attributes in json format

3) for viewing a particular user
make a get request in the url: /users/<id>
# output will be the user and its attributes listed in json format

4) for creating a new contact
provide a http POST request to url: /contacts
json body:
{
  "user_id_1": "<userId1>",
   "user_id_2": "<userId2>"
}  
# timestamp will be added automatically at host using time.Now()

5) for listing all contacts
provide a http GET request to url: /contacts

5) for listing a particular contact with id given
provide a http GET request to url: /contacts?user={id}&infection_timestamp={ts}



