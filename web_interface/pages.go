package web_interface


var requests = `
<!DOCTYPE html>

<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
<script>
  function gen(n) {
    for (var i=0; i<n; i++) {
$.ajax({
  type: 'POST',
  url: "/put/submit?key="+i+"&value="+i,
  aync:false
});
    }
    alert("done")
  }
</script>
<html>
<body>

<input type=text value="0" id="num_requests"/>
<input type=button value ="Generate Requests" onclick="gen(document.getElementById('num_requests').value)" />

</body>
</html>
`

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
<p><a href="#" onClick="onClick()">Click Me</a></p>

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

var bar_graph = `
<!DOCTYPE html>
	<head>
		<script src="http://d3js.org/d3.v3.min.js" charset="utf-8"></script>
		<script src="http://ajax.googleapis.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
	</head>

	<body>
		<style type="text/css">
			#description{
				position: relative;
				left: 200px;
			}

			#description p{
				color: gray;
				font-family: Calibri, Verdana, sans-serif;
				font-size:55px;
				margin: 0;
			}

			.axis path, .axis line {
    				fill: none;
   				stroke: black;
   				shape-rendering: crispEdges;
			}

			.axis text {
    				font-family: sans-serif;
    				font-size: 11px;
			}

			#display{
				margin-left: auto;
				margin-right: auto;
				width: 800px;
			}

			#timer{
				color: gray;
				font-family: Calibri, Verdana, sans-serif;
				font-size:20px;
				margin: 0;
				text-align: center;
			}


		</style>

		<div id="description">
			<p>
				EDHT Partition Size vs. Partition Keys
			</p>
		</div>
		<div id="timer">
			Uptime:
		</div>
		<div id="display">
			<svg id="canvas" width="800" height="500" x ="0" y="0" style="border: 1px solid black">

			</svg>
		</div>
		<script>
			//application parameters
			var SERVER_LINK = "stats/balance/heat";
			var UPDATE_MS_DELAY = 3000;
			var FIRST_DRAW = true;
			var COLD_COLOR = "blue";
			var HOT_COLOR = "red";
			var SCALE_FACTOR = 10000000000000000;

			//application specific state variables
			var date_object = new Date();
			var start_time;

			/*
			* The chart application is the handler for an initial AJAX request, that will
			* either: 1) create the chart elements the first time around, or 2) apply transitions
			* to show new data if it reupdates itself again.
			*
			* Because of the nature of the call, all variables used by the chart must be moved outside
			* the function for persistent state.
			*/

			//chart specific state variables
			var initial_data;
			var partition_data;
			var key_data;

			var canvas;
			var canvas_width;
			var canvas_height;

			var bar_padding = 2;
			var chart_padding = 50;

			var size_max_value;
			var key_max_value;

			//a "scale" for # of partitions; actually stays constant.
			var x_scale;
			//y scale for partition size
			var y_scale;
			//for bar heights
			var height_scale;
			//heat color map
			var color_scale;

			var y_axis;
			//note the x-axis is an ordinal axis for the # of partitions
			var x_axis;

			//d3 selection for the chart bars
			var chart_bar;


			var chart_handler = function(data){
				//update time at the start of each iteration
				var current_time = ((new Date()).getTime() - start_time);
				d3.select("#timer").text("Uptime: " + (current_time/1000).toFixed(2) + "s");

				if(FIRST_DRAW){
					//setup everything for the first time
					initial_data = separate_data(data);
					partition_data = initial_data[1];
					key_data = initial_data[0];

					canvas = d3.select("#canvas");
					canvas_width = canvas.attr("width");
					canvas_height = canvas.attr("height");

					//x discrete, ordinal scale based on # of partitions as given by partition_data.length
					x_scale = d3.scale.ordinal()
							  .domain(d3.range(partition_data.length))
							  .rangeRoundBands([chart_padding, canvas_width - chart_padding], 0.1);

					//map 0 to padded bottom, max partition size to padded top
					size_max_value = d3.max(partition_data);
					y_scale = d3.scale.linear()
							  .domain([0, size_max_value])
							  .range([canvas_height - chart_padding, chart_padding]);

					//map 0 to 0 height, max partition size to amt of plottable height
					height_scale = d3.scale.linear()
							 .domain([0, size_max_value])
							 .range([0, canvas_height - (2 * chart_padding)]);

					key_max_value = d3.max(key_data);
					//map each # of keys in partition from cold color to hot color
					color_scale = d3.scale.linear()
							.domain([0, key_max_value])
							.range([COLD_COLOR, HOT_COLOR]);

					//generate + add the axes given the current data
					y_axis = d3.svg.axis().scale(y_scale).orient("left");
					x_axis = d3.svg.axis().scale(x_scale).orient("bottom");

					canvas.append("g")
						.attr("class", "y_axis axis")
						.attr("transform", "translate(" + (chart_padding - 5) + ", 0)")
						.call(y_axis);

					canvas.append("g")
						.attr("class", "x_axis axis")
						.attr("transform", "translate(0," + (canvas_height - chart_padding) + ")")
						.call(x_axis);

					//add the rectangular bars
					chart_bar = canvas.selectAll("rect").data(partition_data).enter()
							  .append("rect")
							  .attr("id", function(d, i) {return "bar" + i;})
							  .attr("x", function(d, i) {return x_scale(i);})
							  .attr("y", function(d, i) {return y_scale(d);})
							  .attr("width", x_scale.rangeBand())
							  .attr("height", function(d, i){ return height_scale(d);})
							  .style("fill", function(d, i){
								return color_scale(key_data[i]);
							   });

					//since each rectangle has an id equal to its index number in the partition size list,
					//we loop through and just add text labels
					for(var i = 0; i < partition_data.length; i++){
						var bar = d3.select("#bar" + i);

						//must parseInt, because the attributes are strings!
						canvas.append("text")
						.attr("x", parseInt(bar.attr("x")) + parseInt((bar.attr("width") / 2)))
						.attr("y", parseInt(bar.attr("y")) - 5)
						.text(partition_data[i] + "/" + key_data[i])
						.attr("id", "text" + i)
						.style("text-anchor", "middle")
						.style("fill", "red");
					}

					//set a flag so that we will perform updates from now on.
					FIRST_DRAW = false;
				}
				else{
					//
					var new_data = separate_data(data);
					var new_partition_data = new_data[1];
					var new_key_data = new_data[0];

					//update the scales if necessary
					var new_size_max_value = d3.max(new_partition_data);
					var new_key_max_value = d3.max(new_key_data);

					//rescale the y axis and height_scales
					if(new_size_max_value > size_max_value){
						size_max_value = new_size_max_value;
						y_scale.domain([0, new_size_max_value]);
						height_scale.domain([0, new_size_max_value]);
					}

					//similarly, rescale the color scale
					if(new_key_max_value > key_max_value){
						key_max_value = new_key_max_value;
						color_scale.domain([0, new_key_max_value]);
					}

					//recall the axes
					d3.select(".y_axis").transition().duration(1000).call(y_axis);

					//in our construction, we know that the bar for partition i has id "bar" + i
					//and its corresponding text label has id "text" + 1
					for(var i = 0; i < new_partition_data.length; i++){
						//simultaneously?? change the text label to new partition size value
						//and transition rectangle height and y pos
						var new_y_pos = y_scale(new_partition_data[i]);
						var new_height = height_scale(new_partition_data[i]);
						var cur_rect = d3.select("#bar" + i);
						var cur_text = d3.select("#text" + i);

						cur_text.text(new_partition_data[i] + "/" + new_key_data[i]);
						cur_text.transition()
							.attr("y", new_y_pos - 5)
							.duration(1000);

						cur_rect.transition()
							.attr("y", new_y_pos)
							.attr("height", new_height)
							.style("fill", color_scale(new_key_data[i]))
							.duration(1500);
					}

				}


				//issue the next ajax call after waiting UPDATE_MS_DELAY milleseconds
				setTimeout(function() {$.get(SERVER_LINK, chart_handler, "json");}, UPDATE_MS_DELAY);
			}

			//start the recursive function, and initialize initial_data
			//along with initial start time
			start_time = date_object.getTime();
			$.get(SERVER_LINK, chart_handler, "json");

			//------------------------------------------------------------------------------------------

			/*
			* Given data from /stats/balance/heat of form [{keys : ..., partition : ..., }, {same}, ... ]
			* Returns a list of lists [[keys_i],[partition_i]] where ith index of both lists correspond to
			* the value of the attributes for object i.
			*/
			function separate_data(data){
				var size_data = [];
				var key_data = [];
				for(var i = 0; i < data.length; i++){
					size_data.push(Math.floor(data[i]["size"]/SCALE_FACTOR));
					key_data.push(Math.floor(data[i]["keys"]));
				}

				return [key_data, size_data];
			}

		</script>
	</body>
</html>
`
