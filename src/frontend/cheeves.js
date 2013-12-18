$(document).ready(function () {
});


function deliverCheeve(card, badge){


    $('body').append('<div id="cheeve" class="hide-right"></div>');
    $('#cheeve').append('<div id="left"></div>');
    $('#cheeve').append('<div id="right"></div>');
    $('#left').append('<img width="80" height="80" src="image.png" >');
    $('#right').append('<h2>'+badge.Name+'</h2>');

    if(card.Points != badge.PointsRequired){
        $('#right').append('<p>'+card.Points+' of '+badge.PointsRequired+'</p>');
    } else {
        $('#right').append('<p>Edit 4 more features to earn this</p>');
    }
    

    $('#cheeve').click(function(event) {
            $('#cheeve').toggleClass('hide-right');
            console.log("toggle");
    });

    setTimeout(function(){
        $('#cheeve').toggleClass('hide-right');
    },500);

}

        var cheeves = {
                update: function update(url, email, badge) {

                $('#cheeve').remove();
                // url = "http://127.0.0.1:3000";
                // email = "wookoouk@gmail.com";
                // badge = "5";

                var card = getCard(url, email, badge);
                var badge = getBadge(url, badge);

                if(card.Given){

                    alert("You already have this badge");

                } else {deliverCheeve(card, badge);}
            }
        }


function getCard(url, email, badge) {
//        http://127.0.0.1:3000/cards/update/wookoouk@gmail.com/1

var getter = url+"/cards/update/"+email+"/"+badge;

    var card;
    $.ajax({
        url : getter,
        type: "GET",
        dataType:'json',
        async: false,
        success : function(data) {  
            card = {
                Email: data.Email,
                Badge: data.Badge,
                Points: data.Points,
                Given: data.Given
            }
        }
    });
    return card;
}

function getBadge(url, badge) {
//    http://127.0.0.1:3000/badges/1

var getter = url+"/badges/"+badge;

    var badge;
    $.ajax({
        url : getter,
        type: "GET",
        dataType:'json',
        async: false,
        success : function(data) {  
            badge = {
                ID: data.ID,
                Name: data.Name,
                Badge: data.Badge,
                Evidence: data.Evidence,
                PointsRequired: data.PointsRequired
            }
        }
    });
    return badge;
        
}