$.fn.serializeObject = function()
{
   var o = {};
   var a = this.serializeArray();
   $.each(a, function() {
       if (o[this.name]) {
           if (!o[this.name].push) {
               o[this.name] = [o[this.name]];
           }
           o[this.name].push(this.value || '');
       } else {
           o[this.name] = this.value || '';
       }
   });
   return o;
};

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";";
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "empty";
}

function login(e) {
    e.preventDefault();

    let formData = $("#login").serializeObject();
    fetch("/login", {
      method: 'POST',
      credentials: "same-origin",
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    body: JSON.stringify(formData),
    })
    .then(response => { 
      if  (response.status == 200) {
        return response.json();
      }
    })
    .then(data => {
        console.log(data);
        // window.sessionStorage.token = data.token;
        // It's not possible to use the httpOnly option
        // when setting a cookie client side. 
        // setCookie( "authorize_token", data.token, 1 );
        window.location.replace("http://localhost:8080/toodles");
    });

}

function deleteToodle(e, id) {
    e.preventDefault();

    let t = e.target;
    let formData = $(t).closest( "form" ).serializeObject();

    fetch("/toodles/" + id, {
      method: 'DELETE',
      credentials: "same-origin",
      headers: {
      'content-type': 'application/json',
      'Accept': 'application/json'
    },
    body: JSON.stringify(formData),
    })
    .then(data => { 
      if  (data.status == 200) {
        let t = e.target;
        $(t).closest( "li" ).remove();
      }
    });
}

function addToodle(e) {
    e.preventDefault();

    let formData = $("#add-toodle").serializeObject();
    fetch("/toodles", {
      method: 'POST',
      credentials: "same-origin",
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    body: JSON.stringify(formData),
    })
    .then(response => { 
      if  (response.status == 200) {
        return response.json();
      }
    })
    .then(data => {
        let cookie = getCookie("csrf");
        console.log(cookie);
        let toodle = getTootleHTML(data.id, data.title, data.content, cookie);
        $("ul").append(toodle);
    });
}


function updateToodle(e, id) {
    e.preventDefault();

    let t = e.target;
    let formData = $(t).closest( "form" ).serializeObject();

    fetch("/toodles/" + id, {
      method: 'PUT',
      credentials: "same-origin",
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    body: JSON.stringify(formData),
    })
    .then(response => { 
      if  (response.status == 200) {
        return response.json();
      }
    })
    .then(data => {
        console.log($(t).closest( "li" ));
        $(t).closest( "li" ).find(".title").text(data.title);
        $(t).closest( "li" ).find("input[name=title]").val(data.title)
        $(t).closest( "li" ).find("input[name=content]").val(data.content)
    });
}


function getTootleHTML(id, title, content, cookie) {
    let toodle = `
    <li class="list-group-item">

        <a href="/toodles/${id}" onClick="event.stoppropagation();">
            <span class="title" data-toggle="collapse" data-target="#toodleEdit-${id}" aria-expanded="false" aria-controls="toodleEdit">
                ${title}  
            </span>
        </a>
        
        <div class="abs-right">

            <form class="formComplete" action="/toodles/${id}" method="post">
                <input name="csrf" type="hidden" value="${cookie}" >
                <input type="hidden" name="method" value="put">
                <button type="submit" class="icon-check" onClick="complteToodle(event, '${id}')"></button>
            </form>

            <form class="formDelete" action="/toodles/${id}" method="post">
                <input name="csrf" type="hidden" value="${cookie}" >
                <input type="hidden" name="method" value="delete">
                <button type="submit" class="icon-close" onClick="deleteToodle(event, '${id}')"></button>
            </form>

            <label class="expand" data-toggle="collapse" data-target="#toodleEdit-${id}" aria-expanded="false" aria-controls="toodleEdit">
                <input type="checkbox">
                <i class="icon-chevron-down"></i>
            </label>
        </div>

            <div class="collapse" id="toodleEdit-${id}">
                <div class="card card-body">
                    <form action="/toodles/${id}" method="post">
                        <input name="csrf" type="hidden" value="${cookie}" >
                        <input type="hidden" name="method" value="put">
                        <div class="form-group">
                            <input name="title" type="text" value="${title}" class="form-control" id="title" placeholder="${title}">
                        </div>
                        <div class="form-group">
                            <input name="content" type="text" value="${content}" class="form-control" id="content" placeholder="${content}">
                        </div>
                        <button type="submit" class="btn btn-success" onClick="updateToodle(event,'${id}')">Update</button>
                    </form>
                </div>
            </div>
        </li>
    `;

    return toodle;
}