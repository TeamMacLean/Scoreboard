$(document).ready(function () {

    test();

    $('body').append('<div id="cheeve" class="hide-right"></div>');
    $('#cheeve').append('<div id="left"></div>');
    $('#cheeve').append('<div id="right"></div>');
    $('#left').append('<img width="80" height="80" src="image.png" >');
    $('#right').append('<h2>Feature Editor 1/5</h2>');
    $('#right').append('<p>Edit 4 more features to earn this</p>');


    $('#cheeve').click(function (event) {
        $('#cheeve').toggleClass('hide-right');
    });

});

var cheeves = {
    getResponse: function getResponse(url) {
        $.getJSON(url, function (data) {
            console.log(data);

//            {"Email":"wookoouk@gmail.com","Badge":1,"Points":10,"Given":true}







        });
    }
}


function getBadge() {
//    http://127.0.0.1:3000/badges/1
        var badge = {
            ID: data.ID,
            Name: data.Name,
            Badge: data.Badge,
            Evidence: data.Evidence,
            PointsRequired: data.PointsRequired

        }
        }

    function getCard(){
        var card = {
            Email: data.Email,
            Badge: data.Badge,
            Points: data.Points,
            Given: data.Given
        }
    }


    function test() {
        cheeves.getResponse("http://127.0.0.1:3000/cards/update/wookoouk@gmail.com/1");
    }