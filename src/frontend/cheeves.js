$( document ).ready(function() {

	$('body').append('<div id="cheeve" class="hide-right"></div>');
	$('#cheeve').append('<div id="left"></div>');
	$('#cheeve').append('<div id="right"></div>');
	$('#left').append('<img width="80" height="80" src="image.png" >');
	$('#right').append('<h2>Feature Editor 1/5</h2>');
	$('#right').append('<p>Edit 4 more features to earn this</p>');


	$('#cheeve').click(function(event) {
		/* Act on the event */
		$('#cheeve').toggleClass('hide-right');
	});

});