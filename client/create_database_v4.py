#the purpose of this file is to take as input
#a number of courses as described in a file and
#issue a series of put requests to the coordinator
#to create the base database structure

import requests
import json
import re

coordinator_put_link = "http://localhost:14568/put/submit"

#first read each non-commented line of the file and parse
#into a list of dictionaries, each of the following form:
#{course: , course_number , instructor , description , seats_taken , seats_available , room , time}

#match the vertical bars in the string and ignore whitespace; compiled once for performance reasons
#instead of in loop
re_pattern = re.compile("\\s*\\|\\s*")
course_object_list = []
f = open("course_data.txt", "r")
lines = f.readlines()


for line in lines:
	line = line.strip()
	#only process a line that is not empty and has enough chars for a comment
	if len(line) > 2:
		#skip comments
		if line[:2] == "//":
			continue
		
		course_object = {}
		object_data = re_pattern.split(line)

		#must have enough data to fill in the object's attributes
		assert len(object_data) == 8

		#set the corresponding fields
		course_object["course"] = object_data[0]
		course_object["course_number"] = object_data[1]
		course_object["instructor"] = object_data[2]
		course_object["description"] = object_data[3]
		course_object["seats_taken"] = object_data[4]
		course_object["seats_available"] = object_data[5]
		course_object["room"] = object_data[6]
		course_object["time"] = object_data[7]

		#add it to the list of parsed/converted course objects
		course_object_list.append(course_object)

		print json.dumps(course_object)
f.close()

#the first part of the database places the json representation of
#the above in-memory objects with a unique key: i.e. 1, 2, 3 and so on...
#this is so that the javascript client can request the string form, and then
#easily reparse it with d3 afterwards

#the easiest unique key is the index of the object in our parsed object list
for i in range(0, len(course_object_list)):
	#since we use a course_object's index in the course_object_list 
	#as the key in EDHT, we'd like to save that key in the course object so the client
	#can rediscover it and send signup requests!
	course_object_list[i]["db_key"] = str(i)

	#we create an object to hold the key, value pair that will be converted to
	#the "key" and "value" field for Max's web API
	#key is the string form of the unique id, value is the json string representation
	#of a course object
	form_wrapper = {"key": str(i), "value" : json.dumps(course_object_list[i])}
	request_object = requests.get(coordinator_put_link, params=form_wrapper)

	#we then create another entry to represent the lock for updating this particular
	#course; its name is simply the course id + suffix "_lock", with initial value 0
	form_wrapper = {"key": str(i) + "__lock", "value" : 0}
	request_object = requests.get(coordinator_put_link, params=form_wrapper)

#next, we construct "keyword" entries in the database. these are entries where the 
#"key" is a word, and the "value" is the string representation of a json array, 
#i.e. "[0, 5, 2]," which contains the unique key of a course object whose attributes
#include the word in key in its attributes.
#the attributes 

#first compute a set of vocabulary words over the
#course, instructor, description to serve as "keys"; we standardize with lower 
#case strings, and remove punctuation in the description
vocabulary = set()

#a list of lists, where each element is a set consisting of the normalized, cleaned,
#lowercase words of each course object
course_object_word_list = []

#relevant regex split patterns
by_space = re.compile("\\s+")
punctuation_remover = re.compile("[^a-zA-Z\\s]")
for i in range(0, len(course_object_list)):
	course_keyword_list = by_space.split(course_object_list[i]["course"])
	instructor_keyword_list = by_space.split(course_object_list[i]["instructor"])
	#first remove punctuation from the description, and then retrieve the words
	description_keyword_list = by_space.split(re.sub(punctuation_remover, "", course_object_list[i]["description"]))
	course_number = course_object_list[i]["course_number"]

	#add the lowercase of all these words to our vocabulary set, and construct a set to add to the course_object_word_list
	this_object_word_set = set()
	for word in course_keyword_list:
		lowered_word = word.lower()
		vocabulary.add(lowered_word)
		this_object_word_set.add(lowered_word)

	for word in instructor_keyword_list:
		lowered_word = word.lower()
		vocabulary.add(lowered_word)
		this_object_word_set.add(lowered_word)

	for word in description_keyword_list:
		lowered_word = word.lower()
		vocabulary.add(lowered_word)
		this_object_word_set.add(lowered_word)

	#the course number is a word too, and we do the same
	vocabulary.add(course_number)
	this_object_word_set.add(course_number)

	#add our constructed this_object_word_list to the overall course_object_word_list
	course_object_word_list.append(this_object_word_set)

#IMPORTANT: implicit in all these calculations is that we use the index of a loaded course object in the course_object_ist
#as its unique key in the database by construction. and so when we compute which indices a lower case word key maps to,
#the course_object_list must be unchanged.

#first create empty lists for each term in the vocabulary
keyword_group = {}
for term in vocabulary:
	#term is already a string
	keyword_group[term] = []


#we now go through each term in the vocabulary, and we scan through the sets in the course_object_word_list. if a set
#contains this word, then we take the current index and add it to the list of indices the term maps to.
for term in vocabulary:
	i = 0
	for object_set in course_object_word_list:
		if term in object_set:
			keyword_group[term].append(i)
		i += 1

print keyword_group

#after constructing this list, we now send put requests and put records of form
#key: keyword, value: [indices] to complete the keyword mapping entries in the database
for k, v in keyword_group.iteritems():
	form_wrapper = {"key":k , "value": json.dumps(v)}
	request_object = requests.get(coordinator_put_link, params=form_wrapper)

#now, we create student accounts with default empty schedules, i.e. entries of form
#key: student-id, value: [<NO_CLASS>]
for student in ["wfl33"]:
	form_wrapper = {"key":student, "value":json.dumps(["NO_CLASS"])}
	request_object = requests.get(coordinator_put_link, params=form_wrapper)

	#also put locks for the student object
	form_wrapper = {"key":student+"__lock", "value":0}
	request_object = requests.get(coordinator_put_link, params=form_wrapper)


		
