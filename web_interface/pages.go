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

var easter_egg = `
<!DOCTYPE html>
<html>
<body>
<img src="http://i.imgur.com/g07gkAP.jpg" alt="doge">
<audio controls autoplay>
  <source src="https://ia700702.us.archive.org/2/items/gd1977-05-08.ecm33p.moore.berger.miller.117026.flac16/gd77-05-08s2t02.mp3" type="audio/mpeg">
  Your browser does not support the audio element.
</audio>

</body>
</html>
`
