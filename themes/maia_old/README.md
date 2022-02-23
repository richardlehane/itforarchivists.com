# A hugo theme

## Features!
## Two colour schemes

Choose scheme maia or richard with the top level scheme variable. E.g. in config.toml ... 
[params]
scheme = "maia"

## Sidebar shortcode

Use the sidebar shortcode in your content for a two column layout within text. Text within the first sidebar element goes in left column. The text until the closing sidebar element goes in the right column.

E.g. ...

Some full page content. Bla bla.

{{% sidebar "hmmm this will appear in a left column" %}}And this text will appear in a right column.

Nifty!{{% /sidebar %}}

Lastly a bit more full page content.