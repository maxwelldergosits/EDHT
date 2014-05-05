/*
 * db.js implements functions for interacting with the database.
 * before using any of the provided functions, the links to the
 * web API's put and get links must be set with their respective functions.
 * NOTE: Requires jQuery to be included.
*/

var db_put_link="";
var db_get_link="";

/*
 * sets the link to the web API's put page
*/
function db_set_put_link(link){
	db_put_link = link;
}

/*
 * sets the link to the web API's get page
*/
function db_set_get_link(link){
	db_get_link = link;
}

/*
 * db_raw_read takes in a String key to query the database,
 * and returns the exact String value on success, or undefined
 * if the key does not exist.
*/
function db_raw_read(db_key){
	var result = undefined;
	var result_setter = function(data){
		if(data.value != ""){
			result = data;
		}
	}

	$.get(db_get_link, {key: db_key}, result_setter, "json");

	if(result != undefined){
		return result["value"];
	}
	else{
		return undefined;
	}
}

/*
 * db_read takes in a String key to query the database.
 * It returns the parsed JSON data on success, or undefined if the key 
 * does not exist in the database.
*/
function db_read(key){
	var value = db_raw_read(key);
	
	if(value == undefined){
		return undefined;
	}
	else{
		return JSON.parse(value);
	}
}

/*
 * db_write takes in a String key, and a raw String value, and writes exactly
 * these two to the database, with no conversion. 
 * Returns true if the write succeeded, false otherwise.
*/
function db_raw_write(db_key, db_value){
	//response from server is JSON object of form {key: , succ:, value: }
	var response;
	var response_setter = function(data){
		response = data;
	};

	$.get(db_put_link, {key: db_key, value: db_value}, response_setter, "json");

	if(response["succ"] == "true"){
		return true;
	}
	else{
		return false;
	}
}

/*
 * db_write takes in a String key, and a JSON object, and writes that
 * (key, value) pair to the database. 
 * It returns true if the write succeeded, false otherwise.
*/
function db_write(key, value){
	return db_raw_write(key, JSON.stringify(value));
}

/*
 * this version of db_raw_read is only to be used for the implementation
 * of db tas write, and returns the "ov" (consistent old value) of a field
 * behaves exactly like a TAS instruction.
 * This means its return value will either by "1" or "0".
 * NOTE: must only be used on an existing "__lock" field.
*/
function internal_db_tas(db_key, db_value){
	var result;
	
	var result_setter = function(data){
		result = data;
	}

	$.get(db_put_link, {key: db_key, value: db_value, ov:"true"}, result_setter, "json");
	return result["ov"];
}

/*
 * db_tas_write takes in a String key and a JSON value, and attempts to 
 * perform a synchronized write to the key field.
 * Returns 2 if the key exists and the write succeeds,
 * 1 if the key exists and the write failed,
 * 0 if the key does not exist,
 * or -1 if the key exists but does not support synchronized access.
*/
function db_tas_write(key, value){
	//concurrent read is still safe, because i don't do anything
	//with the result and just want to check that the field exists
	var read_result = db_raw_read(key);
	
	if(read_result != undefined){
		//now check for the existence of the field's lock
		var lock_result = db_raw_read(key + "__lock");
		
		if(lock_result != undefined){
			//acquire the TAS lock
			while(internal_db_tas(key + "__lock", "1") == "1"){
			}

			var success_status = db_write(key, value);			

			//make sure we free TAS lock
			while(db_raw_write(key + "__lock", "0") == false){
			}

			if(success_status == true){
				return 2;
			}
			else{
				return 1;
			}
		}
		else{
			return -1;
		}
	}
	else{
		return 0;
	}
}

/*
 * db_acquire_tas_lock takes a name of a key in the database, and will attempt
 * to obtain the "lock" for that field represented by "key" + "__lock". The function
 * will return true when the lock has been held, or false if the key doesn't exist.
*/
function db_acquire_tas_lock(key){
	var lock_name = key + "__lock";
	
	//if either of the two fields don't exist, we can't acquire anything.
	if(db_raw_read(key) == undefined || db_raw_read(lock_name) == undefined){
		return false;
	}

	while(internal_db_tas(lock_name, "1") == "1"){

	}

	return true;
}

/*
 * db_release_tas_lock simply writes a "0" to the "key" + "__lock" field to release
 * the "lock" associated with key. returns true when the lock has been reset, 
 * false if the key doesn't exist.
*/
function db_release_tas_lock(key){
	var lock_name = key + "__lock";

	//since this will be used on lock release, we can be a little looser on failure
	//condition here.
	if(db_raw_read(lock_name) == undefined){
		return false;
	}

	return db_raw_write(lock_name, "0");
}

/*
* db_tas_write_and_create takes in a String key and a JSON value, and 
* will overwrite/create that current (key, value) pair in the database,
* and allocate additional resources to support sychronized writes to 
* that field.
* Returns true on success, false otherwise.
* NOTE: results undefined if you call this function on a key that already
* exists in the table.
*/
function db_tas_write_and_create(key, value){
	//may not be a good idea to implement this; what happens
	//when multiple ppl try to call this on the same key?
}
