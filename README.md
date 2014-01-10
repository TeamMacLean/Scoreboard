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

Examples
======

Ruby
------
class Badge :

```
class Badge
require 'net/http'
require "rubygems"
require "json"


def self.goGet(url, email, badge)

card = getCard(url, email, badge)
badge = getBadge(url, badge)

bundle = card + badge
return bundle
end

def self.getCard(url, email, badge)

getter = "#{url}/cards/update/#{email}/#{badge}";

uri = URI(getter)
req = Net::HTTP.get(uri)
# json = req.to_json
return req.to_s.html_safe
end

def self.getBadge(url, badge)

getter = "#{url}/badges/#{badge}";

uri = URI(getter)
req = Net::HTTP.get(uri)
# json = req.to_json
return req.to_s.html_safe
end

end
```

main.html.erb
```

<script src="/Credit/src/frontend/cheeves.js"></script>
<link href="/Credit/src/frontend/cheeves.css" media="all" rel="stylesheet" type="text/css" />

...

<script>
    <%unless params[:badge].nil?%>
        <% if user_signed_in? %>
            console.log("running cheves");

            <% url = "http://127.0.0.1:3000" %>
            <% email = current_user.email %>
            <% badgep = params[:badge] %>

            card = <%= Badge.getCard(url, email, badgep) %>
            badge = <%= Badge.getBadge(url, badgep) %>
            cheeves.manualUpdate(card, badge);
        <% end %>
    <% end %>
</script>
```
