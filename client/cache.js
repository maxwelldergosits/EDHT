/*
 * cache.js implements user-side caching of the keys and values
 * obtained from the database, with each result cached for a specified amt of time.
*/

var cache_object = {};

/*
 * each entry in the cache_object has the following form: {value: , timestamp: , expiration_duration: }
*/

/*
 * cache_lookup takes a String key as input, and will attempt to lookup its associated value.
 * Returns the value on success, and undefined if: 1) no such (key, value) pair exists or
 * 2) the item has expired (in which case the (key, value) pair will be cleared).
*/
function cache_lookup(key){
	if(cache_object[key] == undefined){
		return undefined;
	}
	else{
		var current_time = (new Date()).getTime();
		var cache_value_object = cache_object[key];
		
		if(current_time - cache_value_object["timestamp"] > cache_value_object["expiration_duration"]){
			//reset the key and return undefined			
			cache_object[key] = undefined;
			return undefined;
		}
		else{
			return cache_value_object["value"];
		}
	}
}

/*
 * cache_add takes a String key, a JSON value from the database, and a duration of time in
 * milleseconds that indicates how long this (key, value) pair is valid. Multiple cache_adds
 * to the same key will result in the latest value associated with the key.
*/
function cache_add(key, value, expiration_time){
	var cache_value_object = {};
	cache_value_object["value"] = value;
	cache_value_object["timestamp"] = (new Date()).getTime();
	cache_value_object["expiration_duration"] = expiration_time;

	cache_object[key] = cache_value_object;
}
