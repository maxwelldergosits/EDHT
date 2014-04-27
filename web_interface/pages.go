package web_interface


var putForm =  `
<!DOCTYPE html>
<html>
<body>

<form action="/put/submit">
Key: <input type="text" name="key" value=""><br>
Value: <input type="text" name="value" value=""><br>
Get Old Value: <input type="checkbox" name="ov" value="true"><br>
<input type="submit" value="Submit">
</form>

<p>Click the "Submit" button and the form-data will be sent to the key value store</p>

</body>
</html>
`

var getForm = `
<!DOCTYPE html>
<html>
<body>

<form action="/get/submit">
Key: <input type="text" name="key" value=""><br>
<input type="submit" value="Submit">
</form>

<p>Click the "Submit" button and the server will retrieve the key specified</p>

</body>
</html>
`


