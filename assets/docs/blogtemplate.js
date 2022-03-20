var title = window.document.getElementById("title");
var timestamp = window.document.getElementById("timestamp");
var sitemap = window.document.getElementById("sitemap");
var main = window.document.getElementById("main");
var contact_form = window.document.getElementById("contact-form");
var content_info = window.document.getElementById("content-info");

var str = main.innerHTML;

function replaceNewlines() {
    // Replace the \n with <br>
    str = str.replace(/(?:\r\n|\r|\n)/g, "<br>");

    // Update the value of paragraph
    main.innerHTML = str;
}