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
                ,
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