API RULES

The application should be able to handle a POST request to create a new animal entry. Request to create a duplicate entry should be
denied
The application should be able to handle a PUT request to update an existing animal or create a new one if the animal in the payload
doesn't exist yet
The application should be able to handle a DELETE request to delete an existing animal. Throw an error if the specified animal doesn't
exist yet
The application should be able to handle a GET request to get a list of all currently existing animals in the system. If no animal is found
then the API should return 404 Not Found
The application should be able to handle a GET request specifically for an animal by using its ID in the path parameter. If the ID
specified doesn’t exist in the storage, the API should return a 404 Not Found
