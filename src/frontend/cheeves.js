$(document).ready(function () {
});


function deliverCheeve(card, badge){

    imageURL = "http://merit.tsl.ac.uk/badges/"

    $('body').append('<div id="cheeve" class="hide-right"></div>');
    $('#cheeve').append('<div id="left"></div>');
    $('#cheeve').append('<div id="right"></div>');

    console.log(imageURL+badge.Badge+".json");

url = imageURL+badge.Badge+".json";
image = null;
$.ajax({
        url : url,
        type: "GET",
        dataType:'json',
        async: false,
        success : function(data) {
            console.log(data);
            image = jQuery.parseJSON( data );
        }
    });

    if(image != null){
        $('#left').append('<img width="80" height="80" src="'+image.image+'" >');
    }
 
    $('#right').append('<h2>'+badge.Name+'</h2>');

    if(card.Points != badge.PointsRequired){
        $('#right').append('<p>'+card.Points+' of '+badge.PointsRequired+'</p>');
    } else {
        $('#right').append('<h3>Badge unlocked!</h3>');
        $('#right').append('<a href='+card.Assert+'>Click here to collect</a>');
    
    }

    
    setTimeout(function(){
        $('#cheeve').toggleClass('hide-right');
    },500);

}

function process(card, badge) {
    $('#cheeve').remove();

                
                if(card.Given && card.Assert.length < 1){
                    console.log(card);
                    console.log("user already has this badge");

                } else {
                    deliverCheeve(card, badge);
                }
}

        var cheeves = {
                manualUpdate: function manualUpdate(card, badge) {
                    process(card, badge);
                },
                update: function update(url, email, badge) {
                    var card = getCard(url, email, badge);
                    var badge = getBadge(url, badge);
                    process(card, badge);
            }
        }


function getCard(url, email, badge) {
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
                Given: data.Given,
                Assert: data.Assert
            }
        }
    });
    return card;
}

function getBadge(url, badge) {
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