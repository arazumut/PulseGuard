$(function() {
    "use strict";  

    $('.knob2').knob({
        'format' : function (value) {
            return value + '';
		 },
	});

	// Return Visitor Information
	$(document).ready(function() {
		var options = {
			chart: {
				height: 300,
				type: 'bar',
			},
			colors: ['#59c4bc', '#637aae'],
			plotOptions: {
				bar: {
					horizontal: false,
					columnWidth: '55%',
					endingShape: 'rounded'	
				},
			},
			dataLabels: {
				enabled: false
			},
			stroke: {
				show: true,
				width: 2,
				colors: ['transparent']
			},
			series: [{
				name: 'Return',
				data: [44, 55, 57, 56, 61, 58, 63, 60, 66]
			}, {
				name: 'New',
				data: [76, 85, 101, 98, 87, 105, 91, 114, 94]
			}],
			xaxis: {
				categories: ['Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'],
			},
			yaxis: {
				title: {
					text: '$ (thousands)'
				}
			},
			fill: {
				opacity: 1

			},
			tooltip: {
				y: {
					formatter: function (val) {
						return "$ " + val + " thousands"
					}
				}
			}
		}

		var chart = new ApexCharts(
			document.querySelector("#Return-Visitor-Information"),
			options
		);

		chart.render();
	});

	// Bounce Rate
	$(document).ready(function() {
		var options = {
			chart: {
				height: 300,
				type: 'line',
				toolbar: {
					show: false,
				},
			},
			colors: ['#59c4bc', '#637aae'],
			series: [{
				name: 'Website Blog',
				type: 'column',
				data: [440, 505, 414, 671, 227, 413, 201, 352, 752, 320, 257, 160]
			}, {
				name: 'Social Media',
				type: 'line',
				data: [23, 42, 35, 27, 43, 22, 17, 31, 22, 22, 12, 16]
			}],
			stroke: {
				width: [0, 4]
			},        
			// labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
			labels: ['01 Jan 2001', '02 Jan 2001', '03 Jan 2001', '04 Jan 2001', '05 Jan 2001', '06 Jan 2001', '07 Jan 2001', '08 Jan 2001', '09 Jan 2001', '10 Jan 2001', '11 Jan 2001', '12 Jan 2001'],
			xaxis: {
				type: 'datetime'
			},
			yaxis: [{
				title: {
					text: 'Website Blog',
				},

			}, {
				opposite: true,
				title: {
					text: 'Social Media'
				}
			}]
		}
		var chart = new ApexCharts(
			document.querySelector("#Bounce-Rate"),
			options
		);

		chart.render();
	});

	// Map
	var mapData = {
		"US": 298,
		"SA": 200,
		"AU": 760,
		"IN": 2000000,
		"GB": 120,
	};
	if( $('#world-map-markers').length > 0 ){
		$('#world-map-markers').vectorMap(
		{
			map: 'world_mill_en',
			backgroundColor: 'transparent',
			borderColor: '#fff',
			borderOpacity: 0.25,
			borderWidth: 0,
			color: '#e6e6e6',
			regionStyle : {
				initial : {
				fill : '#eaeaea'
				}
			},

			markerStyle: {
			initial: {
							r: 5,
							'fill': '#fff',
							'fill-opacity':1,
							'stroke': '#000',
							'stroke-width' : 1,
							'stroke-opacity': 0.4
						},
				},
		
			markers : [{
				latLng : [21.00, 78.00],
				name : 'INDIA : 350'
			
			},
			{
				latLng : [-33.00, 151.00],
				name : 'Australia : 250'
				
			},
			{
				latLng : [36.77, -119.41],
				name : 'USA : 250'
			
			},
			{
				latLng : [55.37, -3.41],
				name : 'UK   : 250'
			
			},
			{
				latLng : [25.20, 55.27],
				name : 'UAE : 250'
			
			}],

			series: {
				regions: [{
					values: {
						"US": '#49c5b6',
						"SA": '#667add',
						"AU": '#50d38a',
						"IN": '#60bafd',
						"GB": '#ff758e',
					},
					attribute: 'fill'
				}]
			},
			hoverOpacity: null,
			normalizeFunction: 'linear',
			zoomOnScroll: false,
			scaleColors: ['#000000', '#000000'],
			selectedColor: '#000000',
			selectedRegions: [],
			enableZoom: false,
			hoverColor: '#fff',
		});
	}
});