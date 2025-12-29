$(function() {
    "use strict";

    // notification popup
    toastr.options.closeButton = true;
    toastr.options.positionClass = 'toast-bottom-right';
    toastr.options.showDuration = 1000;
    toastr['info']('Hello, welcome to Iconic, a unique admin Template.');

    // Google Analytics Dashboard
    $(document).ready(function(){
        var chart = c3.generate({
            bindto: '#Google-Analytics-Dashboard', // id of chart wrapper
            data: {
                columns: [
                    // each columns data
                    ['data1', 11, 8, 22, 18, 19, 6, 17, 11, 17, 32, 9, 12],
                    ['data2', 7, 7, 5, 7, 9, 12, 8, 22, 18, 19, 6, 17]
                ],
                type: 'bar', // default type of chart
                colors: {
                    'data1': Iconic.colors["theme-cyan1"],
                    'data2': Iconic.colors["theme-cyan2"]
                },
                names: {
                    // name of each serie
                    'data1': '2019',
                    'data2': '2020'
                }
            },
            axis: {
                x: {
                    type: 'category',
                    // name of each category
                    categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'July', 'Aug', 'Sept', 'Oct', 'Nov', 'Dec']
                },
            },
            bar: {
                width: 20
            },
            legend: {
                show: false, //hide legend
            },
            padding: {
                left: 0,
                right: 0,
                bottom: 0,
                top: 0
            },
        });
    });

    //  Use by Device
    $(document).ready(function(){
        var chart = c3.generate({
            bindto: '#Use-by-Device', // id of chart wrapper
            data: {
                columns: [
                    // each columns data
                    ['data1', 50],
                    ['data2', 35],
                    ['data3', 15],
                ],
                type: 'donut', // default type of chart
                colors: {
                    'data1': Iconic.colors["theme-cyan1"],
                    'data2': Iconic.colors["theme-cyan2"],
                    'data3': Iconic.colors["theme-cyan3"]
                },
                names: {
                    // name of each serie
                    'data1': 'Desktop',
                    'data2': 'Mobile',
                    'data3': 'Tablet',
                }
            },
            axis: {
            },
            legend: {
                show: true, //hide legend
            },
            padding: {
                bottom: 0,
                top: 0
            },
        });
    });

    // Use by Audience
    $(document).ready(function(){
        var chart = c3.generate({
            bindto: '#Use-by-Audience', // id of chart wrapper
            data: {
                columns: [
                    // each columns data
                    ['data1', 55],
                    ['data2', 30],
                    ['data3', 20],
                ],
                type: 'donut', // default type of chart
                colors: {
                    'data1': Iconic.colors["theme-purple1"],
                    'data2': Iconic.colors["theme-purple2"],
                    'data3': Iconic.colors["theme-purple3"]
                },
                names: {
                    // name of each serie
                    'data1': 'Male',
                    'data2': 'Female',
                    'data3': 'Other',
                }
            },
            axis: {
            },
            legend: {
                show: true, //hide legend
            },
            padding: {
                bottom: 0,
                top: 0
            },
        });
    });
});