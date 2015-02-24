var init = function(update, onInit) {
    var prmstr = window.location.search.substr(1)
    var prmarr = prmstr.split("&");
    var params = {};

    for ( var i = 0; i < prmarr.length; i++) {
        var tmparr = prmarr[i].split("=");
        params[tmparr[0]] = tmparr[1];
    }

    var refresh = function() {
        var bounds = map.getBounds();

        var ne = bounds.getNorthEast();
        var sw = bounds.getSouthWest();

        var url = '/points?topLeftLat=' + ne.lat() + '&topLeftLon=' + sw.lng() + '&bottomRightLat=' + sw.lat() + '&bottomRightLon=' + ne.lng() + '&index=' + params["index"];

        console.log(url);

        $.getJSON(url,
            function(data) {
                update(data);
             }
        );
    }

    function initialize() {
        var mapOptions = {
            zoom: 12,
            center: new google.maps.LatLng(51.508742,-0.118318),
            mapTypeId: google.maps.MapTypeId.ROADMAP
        };
        map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);
        google.maps.event.addListener(map, 'idle', refresh);

        if (onInit) {
            onInit(map, params)
        }
    }

    if (params["refresh"]) {
        var refreshCycle = function() {
            refresh();
            setTimeout(refreshCycle, 1000);
        }
        setTimeout(refreshCycle, 2000);
    }

    $(document).ready(initialize);
}