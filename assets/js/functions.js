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


function deleteToodle(e, id) {
    e.preventDefault();

    fetch("/todos/" + id, {
      method: 'DELETE',
      headers: {
      'content-type': 'application/json',
      'Accept': 'application/json'
    },
    })
    .then(data => { 
      if  (data.status == 200) {
        var t = e.target;
        $(t).closest( "li" ).remove();
      }
    });
}

function addToodle(e) {
    e.preventDefault();

    let data = $("#add-toodle").serializeObject();
    console.log(data);
    fetch("/todos", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    body: JSON.stringify(data),
    })
    .then(response => { 
      if  (response.status == 200) {
        return response.json();
      }
    })
    .then(data => {
        console.log(data);
        console.log(data.id);
        let toodle = getTootleHTML(data.id, data.title, data.note);
        $("ul").append(toodle);
    });
}



function getTootleHTML(id, title, note) {
    let toodle = `
    <li class="list-group-item">

        ${title} 
        
        <div class="abs-right">

            <form class="formDelete" action="/todos/${id}" method="post">
            <input type="hidden" name="method" value="delete">
            <button type="submit" class="icon-close" onClick="deleteToodle(event, ${id})"></button>
            </form>

            <label class="expand" data-toggle="collapse" data-target="#toodleEdit-${id}" aria-expanded="false" aria-controls="toodleEdit" onClick="() => e.preventDefault()">
                <input type="checkbox">
                <i class="icon-chevron-down"></i>
            </label>
        </div>

            <div class="collapse" id="toodleEdit-${id}">
                <div class="card card-body">
                    <form action="/todos/${id}" method="post">
                    <input type="hidden" name="method" value="put">
                    <div class="form-group">
                        <input name="title" type="text" value="${title}" class="form-control" id="title" placeholder="${title}">
                    </div>
                    <div class="form-group">
                        <input name="note" type="text" value="${note}" class="form-control" id="note" placeholder="${note}">
                    </div>
                    <button type="submit" class="btn btn-success">Update</button>
                </form>
                </div>
            </div>
        </li>
    `;

    return toodle;
}