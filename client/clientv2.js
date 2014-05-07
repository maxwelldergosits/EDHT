/*
 * client.js implements the code the controls the sclient html file.
 * NOTE: Requires d3 to be loaded.
*/

var logged_in = false;
var user_id = undefined;

/*
 * initializes the client, cache, and database with the correct links to the coordinator.
*/
function client_setup(put_link, get_link){
	db_set_put_link(put_link);
	db_set_get_link(get_link);

	$.ajaxSetup({ async: false});
	console.log("client init");
}

/*
 * when called, takes given ID in the "login" textbox, and sees if its
 * an existing user in the database. If so, it will set the logged_in flag
 * to true, and reveal + enable the functional parts of the client.
*/
function client_login(){
	var id_input = document.getElementById("login").value.trim();
	if(db_read(id_input) != undefined){
		user_id = id_input;
		logged_in = true;
		d3.select(".logged_in").style("display", "block");	
	}
	else{
		alert("Could not log in with username " + id_input);
	}
}

/*
 * when called, takes the logged in student id, queries the database
 * for the student's course list, then separately queries the database
 * for each course, caches the data, and then displays the results.
 * timeout indicates how long we want to cache the user's schedule.
*/

//the client's course list is either ["NO_CLASS"] or a list of id's
function client_display_schedule(first_check_cache, cache_timeout){
	if(logged_in == true){
		var user_course_list;
		var course_object_list = [];

		if(first_check_cache == true){
			user_course_list = cache_lookup(user_id);
		}

		//force refresh or if we don't check cache
		if(!first_check_cache || (user_course_list == undefined)){
			user_course_list = db_read(user_id);
		}


		if(user_course_list != undefined){

			//alert the user if no courses to display
			if(user_course_list.length == 1 && user_course_list[0] == "NO_CLASS"){
				alert("You are currently not signed up for any courses!");
				//erase anything that might be left from previous operations
				d3.selectAll(".schedule_table").data([]).exit().remove();
			}
			else{
				//everytime we check the cache or do a read, we get 
				//a course object back;
				//we just need to add a num attribute for my display code
				for(var i = 0; i < user_course_list.length; i++){
					var course_object;
					if(first_check_cache == true){
						//remember the key must be a string
						course_object = cache_lookup(user_course_list[i]);
					}

					if(!first_check_cache || course_object == undefined){
						course_object = db_read(user_course_list[i]);
					}

					//if we haven't failed, we have a course object
					if(course_object != undefined){
						course_object["num"] = i+1;
						course_object_list.push(course_object);
					}
				}

				//we can now cache the user's course list!
				cache_add(user_id, user_course_list, cache_timeout);

				//remove the previous schedule table, and draw the new one
				d3.selectAll(".schedule_table").data([]).exit().remove();
				create_table(user_course_list, "schedule", course_object_list, ["course", "instructor", "seats_available", "seats_taken", "time"], "schedule_table", false, true);

			}
			

		}
		else{
			alert("Could not load your schedule!");
		}
		
	}
}

function sanitize_input(str){
	var trimmed_str = str.trim();
	var no_punctuation = trimmed_str.replace(/[^a-zA-Z0-9\s+]/g, "");
	var all_lowered = no_punctuation.toLowerCase();

	return all_lowered;
}

/*
 * when called, client_search reads in the search options and keywords,
 * searches in lowercase (normalized form for the db)
 * and either queries the cache or the database for each course, and displays
 * their results.
*/
function client_search(first_check_cache, cache_timeout){
	if(logged_in == true){
		//split according to spaces, remove punctuation, and search
		//the lower case forms of the words in the database/cache

		var search_keyword_string = sanitize_input(document.getElementById("keyword").value);
		var search_number_string = sanitize_input(document.getElementById("number").value);

		var keyword_list = search_keyword_string.split(/\s+/);
		var number_list = search_number_string.split(/\s+/);
		var unclean_terms_list = keyword_list.concat(number_list);
		var terms_list = [];

		//need to filter out empty strings from unclean_terms_list
		for(var i = 0; i < unclean_terms_list.length; i++){
			if(unclean_terms_list[i] != ""){
				terms_list.push(unclean_terms_list[i]);
			}
		}

		console.log(terms_list);

		//once we have a list of terms to search for, check cache or consult db
		//to find the course id's of the results
		var course_id_list = [];
		//hashset to make sure we don't add duplicate course id's to the list
		var course_seen = {};

		for(var i = 0; i < terms_list.length; i++){
			var result_list = [];

			if(first_check_cache){
				result_list = cache_lookup(terms_list[i]);
			}

			if(!first_check_cache || (result_list == undefined)){
				result_list = db_read(terms_list[i]);
			}


			//only add the ids in result if it's not undefined and 
			//we haven't seen it yet
			console.log("res");
			console.log(result_list);
			if(result_list != undefined){
				//cache the result_list for this term if it's not undefined
				cache_add(terms_list[i], result_list, cache_timeout);


				for(var j = 0; j < result_list.length; j++){
					if(course_seen[result_list[j]] == undefined){
						console.log("Push: " + result_list[j]);
						course_id_list.push(result_list[j]);
						//mark this id as seen
						course_seen[result_list[j]] = 1;
					}
				}
			}

		}

		console.log("crs");
		console.log(course_id_list);

		var course_object_list = [];
		//now we either consult the cache to obtain information on each course_id
		//and recache the obtained information with timeout cache_timeout
		//NOTE: we DO NOT reset the timeout duration on previously cached elements
		for(var i = 0; i < course_id_list.length; i++){
			var from_cache = true;
			var course_object;

			if(first_check_cache){
				course_object = cache_lookup(course_id_list[i]);
			}


			//if we enter this branch to set the course_object, this means
			//the value is not from the cache
			if(!first_check_cache || course_object == undefined){
				from_cache = false;
				course_object = db_read(course_id_list[i]);
			}

			console.log("crs" + course_object);

			//if the object is legit, add it to the cache with cache_timeout duration
			//and append to our course_object_list + add the numerical identifier
			//need to draw our table.
			if(course_object != undefined){
				if(!from_cache){
					cache_add(course_id_list[i], course_object, cache_timeout);
				}
				course_object["num"] = i+1;
				course_object_list.push(course_object);
			}
		}

		var user_course_id_list = db_read(user_id);
		//draw out the table if the course_object_list is well defined
		if(course_object_list.length >= 1){
			d3.selectAll(".search_table").data([]).exit().remove();
			create_table(user_course_id_list, "result", course_object_list, ["course", "course_number", "instructor", "description", "seats_taken", "seats_available", "room", "time"], "search_table", true, false);
		}

	}
}

//course_id_list used to control the fact that we can't enroll again after being signed up for a class already.
function create_table(course_id_list, dest_id, data, headers, custom_class, add_signup_buttons, add_remove_buttons){
	var data_table = d3.select("#"+dest_id).append("table")
				.attr("class", custom_class);

	//create a special header row
	var header_row = data_table.append("tr")
					.attr("class", custom_class + " row_header");

	for(var i = 0; i < headers.length; i++){
		header_row.append("th").attr("class", custom_class + " header_item")
					.style("background-color", "#FF0000")
					.style("color", "#FFFFFF")
					.text(headers[i]);
	}

	//variables so that we don't keep adding enroll and remove headers
	var added_enroll_heading = false;
	var added_remove_heading = false;

	//create a row for each object in the list, and display attributes specified
	//by header_list in the cell
	for(var i = 0; i < data.length; i++){
		var data_row = data_table.append("tr")
					.attr("class", custom_class);
		
		//add a data cell for each attribute in each row
		for(var j = 0; j < headers.length; j++){
			data_row.append("td").attr("class", custom_class + " entry").text(data[i][headers[j]]);
		}

		//if we choose to add signup buttons, we add a signup button whose action is the signup function
		//we can't add an enroll button to a course we've signed up for already
		if(data[i]["seats_available"] > 0 && add_signup_buttons && course_id_list.indexOf(parseInt(data[i]["db_key"])) == -1){
			console.log("cid list " + course_id_list);			
			console.log("indexOf " + data[i]["db_key"]);

			var button_cell = data_row.append("td");

			if(!added_enroll_heading){
				header_row.append("th").attr("class", custom_class + " header_item")
					.style("background-color", "#FF0000")
					.style("color", "#FFFFFF")
					.text("enroll");

				added_enroll_heading = true;

			}
			button_cell.append("input")
					.attr("type", "button")
					.attr("class", custom_class + " add_" + data[i]["db_key"])
					.attr("onclick", "minus_course_seat(user_id," + data[i]["db_key"] + ",1)");
		}	

		if(add_remove_buttons){
			var button_cell = data_row.append("td");

			if(!added_remove_heading){
				header_row.append("th").attr("class", custom_class + " header_item")
					.style("background-color", "#FF0000")
					.style("color", "#FFFFFF")
					.text("remove");

				added_remove_heading = true;

			}
			button_cell.append("input")
					.attr("type", "button")
					.attr("class", custom_class + " remove_" + data[i]["db_key"])
					.attr("onclick", "minus_course_seat(user_id," + data[i]["db_key"] + "," + (-1) + ")");
		}

		
	}

}

/*given the course_id, signs up for the course with tas capabilities.
 *adds the course_id representing a certain course to the user's 
 *schedule_list
*/
function signup_course_id(user, course_id){
	if(logged_in == true){
		//this flag set to true when we have taken a seat!
		var signup_success = false;

		//acquire the lock on the course_id field
		if(db_acquire_tas_lock(course_id) == false){
			//in theory this should never be run, except when we try
			//to erroneously get a lock on a nonexistent key
			alert("TAS acquire on nonexistent key!");
		}

		//now we must be consistent and check that there are actually seats!
		var course_object = db_read(course_id);

		if(course_object != undefined){
			if(course_object["seats_available"] > 0){
				course_object["seats_available"] -= 1;
				course_object["seats_taken"] += 1;
				
				//update the object!
				var write_success = db_write(course_id, course_object);
				if(write_success){
					//we have successfully reserved a spot
					signup_success = true;
				}
				else{
					alert("There was an error in reserving a spot!");	
				}
			}
		}
		else{
			alert("There was an error in signing up for the course.");
		}		

		//release the lock
		if(db_release_tas_lock(course_id) == false){
			//just some error checking
			alert("TAS release failed on key!");
		}

		if(signup_success){
			//i reason that it is safe to perform unsychronized write to the user's data
			//the whole point is that it's fine to be inconsistent on schedule data until
			//we need to signup for a course
			var user_course_list = db_read(user);
			user_course_list.push(course_id);
			db_write(user, user_course_list);
			user_course_list = db_read(user);

			//should handle errors with writing to user's schedule
			//refresh the user's schedule; force a refresh.
			client_display_schedule(false, 5000);
		}
		else{
			alert("Could not update your schedule!");
		}
	}		
}

/*given the course_id, changes the course's seats with tas capabilities.
 *if value is positive, we subtract seats from the course
 *if value is negative, we add course seats.
*/
function minus_course_seat(user, course_id, value){
	value = parseInt(value);
	if(logged_in == true){
		//this flag set to true when we have taken a seat!
		var signup_success = false;

		//acquire the lock on the course_id field
		if(db_acquire_tas_lock(course_id) == false){
			//in theory this should never be run, except when we try
			//to erroneously get a lock on a nonexistent key
			alert("TAS acquire on nonexistent key!");
		}

		//now we must be consistent and check that there are actually seats!
		var course_object = db_read(course_id);

		if(course_object != undefined){
			if(course_object["seats_available"] > 0){
				//must be careful with types here.
				course_object["seats_available"] = parseInt(course_object["seats_available"]) - value;
				course_object["seats_taken"] = parseInt(course_object["seats_taken"]) + value;
				
				//update the object!
				var write_success = db_write(course_id, course_object);
				if(write_success){
					//we have successfully reserved a spot
					signup_success = true;
				}
				else{
					alert("There was an error in reserving a spot!");	
				}
			}
			else{
				alert("The course is full!");
			}
		}
		else{
			alert("There was an error in signing up for the course.");
		}		

		//release the lock
		if(db_release_tas_lock(course_id) == false){
			//just some error checking
			alert("TAS release failed on key!");
		}

		if(signup_success){
			//i reason that it is safe to perform unsychronized write to the user's data
			//the whole point is that it's fine to be inconsistent on schedule data until
			//we need to signup for a course
			var user_course_list = db_read(user);

			//if we take a course seat away, we give it to the user
			if(value > 0){
				user_course_list.push(course_id);
				db_write(user, user_course_list);
			}
			else{
				//else we remove the course from the user's schedule
				var new_user_course_list = [];
				for(var i = 0; i < user_course_list.length; i++){
					if(user_course_list[i] != course_id){
						new_user_course_list.push(user_course_list[i]);
					}
				}
				db_write(user, new_user_course_list);
			}
			
			
			user_course_list = db_read(user);

			//hide the button so the user can't remove or reenroll again
			if(value > 0){
				d3.selectAll(".add_"+course_id).data([]).exit().remove();
			}
			else{
				d3.selectAll(".remove_"+course_id).data([]).exit().remove();
			}

			//should handle errors with writing to user's schedule
			//refresh the user's schedule; force a refresh.
			client_display_schedule(false, 5000);
		}
		else{
			alert("Could not update your schedule!");
		}
	}		
}


