$.fn.serializeObject = function()
{
   var o = {};
   var a = this.serializeArray();
   $.each(a, function() {
       if (o[this.name]) {
           if (!o[this.name].push) {
               o[this.name] = [o[this.name]];
           }
           o[this.name].push(this.value.trim() || '');
       } else {
           o[this.name] = this.value.trim() || '';
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

    if(!validateLoginForm()) {
        return;
    }

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
        if  ([200,302,401].includes(response.status)) {
            return response.json();
        } else {
            handleError(response.status)
        }
    })
    .then(data => {
        if (data.error) {
            setError(data.error);
        } else {
            window.location.replace(location.origin + "/toodles");
        }
        // window.sessionStorage.token = data.token;
        // It's not possible to use the httpOnly option
        // when setting a cookie client side. 
        // setCookie( "authorize_token", data.token, 1 );
    });

}

function signup(e) {
    e.preventDefault();

    if(!validateSignupForm()) {
        return;
    }

    let formData = $("#signup").serializeObject();
    fetch("/signup", {
      method: 'POST',
      credentials: "same-origin",
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    body: JSON.stringify(formData),
    })
    .then(response => { 
        if  ([201,302,400].includes(response.status)) {
            return response.json();
        } else {
            handleError(response.status)
        }
    })
    .then(data => {
        if (data.error) {
            setError(data.error)
        } else {
            $(".card-body").addClass("text-center");
            $("#signup").replaceWith("Your account has been created. Please <a href=\"login\">login</a>.")
        }
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
    .then(response => { 
        if  (response.status == 200) {
            return response.json();
        } else {
            handleError(response.status)
        }
    })
    .then(data => {
        if (data.error) {
            setError(error);
        } else {
            let t = e.target;
            $(t).closest( "li" ).remove();
            decriment('.activeToodles', 1);
        }
    });
}

function addToodle(e) {
    e.preventDefault();
    let t = e.target;
    let form = $(t).closest("form");
    if(!validateToodle(form)) {
        return;
    }

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
        if  ([200,400].includes(response.status)) {
            return response.json();
        } else {
            handleError(response.status)
        }
    })
    .then(data => {
        if (data.error) {
            setError(data.error);
        } else {
            data = data.payload;
            resetForm(form);
            let cookie = getCookie("csrf");
            let toodle = getToodleHTML(data.id, data.title, data.content, cookie);
            $("ul").append(toodle);
            increment('.activeToodles', 1);
        }
    });
}


function updateToodle(e, id) {
    e.preventDefault();

    let t = e.target;
    let form = $(t).closest("form");
    if(!validateToodle(form)) {
        return;
    }

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
        if  (response.status == 201) {
            return response.json();
        } else {
            handleError(response.status)
        }
    })
    .then(data => {
        if (data.error) {
            setError(data.error);
        } else {
            data = data.payload;
            $(t).closest( "li" ).find(".title").text(data.title);
            $(t).closest( "li" ).find("input[name=title]").val(data.title)
            $(t).closest( "li" ).find("input[name=content]").val(data.content)
        }
    });
}

function completeToodle(e, id) {
    e.preventDefault();

    let t = e.target;
    let formData = $(t).closest( "form" ).serializeObject();

    fetch("/toodles/" + id + "/complete", {
      method: 'PUT',
      credentials: "same-origin",
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    body: JSON.stringify(formData),
    })
    .then(response => { 
        if  (response.status == 201) {
            return response.json();
        } else {
            handleError(response.status)
        }
    })
    .then(data => {
        if (data.error) {
            setError(data.error);
        } else {
            let toodle = $(t).closest( "li" );
            toodle.css('background', 'green')
            toodle.fadeOut();
            decriment('.activeToodles', 1);
            increment('.completedToodles', 1);
        }
    });
}

function getToodleHTML(id, title, content, cookie) {
    let toodle = `
    <li class="list-group-item">

        <a href="/toodles/${id}" onClick="checkInput(event);">
            <span class="title" data-toggle="collapse" data-target="#toodleEdit-${id}" aria-expanded="false" aria-controls="toodleEdit">
                ${title}  
            </span>
        </a>
        
        <div class="abs-right">

            <form class="formComplete" action="/toodles/${id}" method="post">
                <input name="csrf" type="hidden" value="${cookie}" >
                <input type="hidden" name="method" value="put">
                <button type="submit" class="icon-check" onClick="completeToodle(event, '${id}')"></button>
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
                            <input name="title" type="text" value="${title}" class="form-control" id="title" placeholder="${title}" required>
                            <small class="form-text text-danger titleHelp"></small>
                        </div>
                        <div class="form-group">
                            <textarea name="content" type="text" class="form-control" id="content" placeholder="${content}" rows="3">${content}</textarea>
                        </div>
                        <button type="submit" class="btn btn-success" onClick="updateToodle(event,'${id}')">Update</button>
                    </form>
                </div>
            </div>
        </li>
    `;

    return toodle;
}

function checkInput(event) {
    event.preventDefault();
    var t = event.target;
    var currentState = $(t).closest('li').find("input[type=checkbox]").prop("checked");
    $(t).closest('li').find("input[type=checkbox]").prop("checked", !currentState);
}

function handleError(status) {
    if (status == 401) {
        window.location.replace(location.origin + "/login");
    } else {
        console.log(error)
    }
}

function setError(err) {
    for (let prop in err) {
        let field =  "#" + prop.toLowerCase() + "Help";
        $(field).text(err[prop]);
    }
}

function decriment(target, n) {
    let cur = $(target).find("span").text();
    cur = parseInt(cur, 10);
    $(target).find("span").text(cur - n);
}

function increment(target, n) {
    let cur = $(target).find("span").text();
    cur = parseInt(cur, 10);
    $(target).find("span").text(cur + n);
}

function validateEmail(target) {
    let email = $(target).val();
    return email.includes('@');
  }

function validatePassword(target) {
    let password = $(target).val();
    return password.length > 3;
}

function validateTitle(target) {
    return $(target).find('input[name="title"]').val().trim() !== "";
}

function validateLoginForm() {

    let validEmail = validateEmail('#email');
    let validPassword = validatePassword('#password');

    validEmail ? $("#emailHelp").text("") : $("#emailHelp").text("Please enter a valid email.");
    validPassword ? $("#passwordHelp").text("") : $("#passwordHelp").text("Don't forget your password.");

    return validEmail && validPassword;
}

function validateSignupForm() {

    let validEmail = validateEmail('#email');
    let validPassword = validatePassword('#password');

    validEmail ? $("#emailHelp").text("") : $("#emailHelp").text("Please enter a valid email.");
    validPassword ? $("#passwordHelp").text("") : $("#passwordHelp").text("Your password should be at least 4 characters.");

    return validEmail && validPassword;
}

function validateToodle(form) {
    let validTitle = validateTitle(form);

    validTitle ? $(form).find(".titleHelp").text("") : $(form).find(".titleHelp").text("Title should not be  empty.");

    return validTitle;
}

function resetForm(form) {
    $(form).find('[name="title"]').val("");
    $(form).find('[name="content"]').val("");
}