var AnalyticsID="AID-12345678";
var ua = window.navigator.userAgent;
var url = window.location.href;
var referrer = document.referrer

var d = new Date,
    dformat = [d.getFullYear(),
               d.getMonth()+1,
               d.getDate()].join('-')+' '+
              [d.getHours(),
               d.getMinutes(),
               d.getSeconds()].join(':');

// send to analytics server
var sendUrl = encodeURI("http://192.168.99.106:31776/analytics.gif?ID=" + AnalyticsID + "&UserAgent=" + ua + "&URL=" + url + "&Referrer=" + referrer + "&Time=" + dformat);

// add to image src
var obj = document.getElementById("analytics");
obj.src = sendUrl;
