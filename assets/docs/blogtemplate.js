async function getUsers() {
    let url = 'post1.json';
    try {
        let res = await fetch(url);
        return await res.json();
    } catch (error) {
        console.log(error);
    }
}


var title = window.document.getElementById("title");
var timestamp = window.document.getElementById("timestamp");
var sitemap = window.document.getElementById("sitemap");
var main = window.document.getElementById("main");
var contact_form = window.document.getElementById("contact-form");
var content_info = window.document.getElementById("content-info");

document.body.addEventListener("click", function() {
    title.innerText = getUsers().title;
    timestamp.textContent = String(Date.now());
    main.textContent = getUsers().main;
});