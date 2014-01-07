Scoreboard
======

Scoreboard tracks the effort that has been contributed by a user as points, when the user has collected enough points towards any particular badge it requests a badge from Merit (https://github.com/wookoouk/merit) and displays it to the user for collection.

Scoreboard consists of two parts, the front end and the back end.
The backend is a RESTful web service that accepts POST and GET requests.

Backend
======

When a message is received on the URL '/cards/update/:email/:badge' Scoreboard will look for a record (or 'card') with a matching email address and badge id, if one is not found one will be created and one point will be given to it, if a card does exist it will add a point, if the additional point brings the score to the requirement to award the badge it will request a new badge from Merit and return its url in the response.

Frontend ('Cheeves')
======

Comming soon...

API Ref
======

Comming soon...
