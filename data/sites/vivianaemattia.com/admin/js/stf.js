var attendeesConfirmed = 0;
var attendeesNotConfirmed = 0;
var attendeesWaiting = 0;
var attendeesTotal = 0;
var attendeesSent = 0;
var attendeesToSend = 0;
$(document).ready(function() {

    $('#import').click(function () {
        $('#loading').fadeIn(200, function () {

            var formData = new FormData();
            formData.append("action", "import");

            // Add files
            var file = document.getElementById('importFile').files[0];
            formData.append("file", file);


            $.ajax({
                url: 'import.html',
                data: formData,
                type: 'POST',
                contentType: false,
                processData: false,
                dataType: 'html',
                success: function (resp) {
                    if (resp == 'ok') {
                        $('#loading').fadeOut(200, function () {
                            toastr.success('DB uploaded successfully');
                            $('#importModalCenter').trigger('reset');
                            $('#importModalCenter').modal('toggle');
                        });
                    }else {
                        $('#loading').fadeOut(200, function () {
                            toastr.warning(resp);
                        });
                    }
                },
                error: function () {
                    $('#loading').fadeOut(200, function () {
                        toastr.error('Connection Error');
                    });
                },
                complete: function(){
                    document.getElementById('importFile').value  = "";
                }
            });

        });
    });

    $('#export').click(function () {
        location.href = "export.html";
    });



    $('#buttonViewItems').click(function(){
  location.href = "attendees.html";
});
    $('#buttonViewFamilies').click(function(){
        location.href = "families.html";
    });
});
function resetModalInsert(){
  $('#insertModal').trigger('reset');
  $('#insertModal').modal('toggle');
  $('#inputName').val("");
  $('#inputCount').val("1")
  $('#inputExpiration').val("")
}

function updateItem(id){
    $('#loading').fadeIn(200, function(){
        var parameters = {
            'status' : $("#inputStatus"+id).val(),
            'id' : parseInt(id, 10)
        }
        parameters = JSON.stringify(parameters);
        $.ajax({
            url: "/admin/attendees.html",
            data: parameters,
            type: 'POST',
            cache: false,
            contentType: 'application/x-www-form-urlencoded; charset=UTF-8',
            dataType: 'html',
            success: function(response){
                if(response == 'Updated'){
                    $('#loading').fadeOut(200,function(){
                        toastr.success('Successfully updated');
                        location.reload();
                    });
                }else{
                    $('#loading').fadeOut(200,function(){
                        toastr.warning('Warning: '+ response);
                    });
                }
            },
            error: function(){
                $('#loading').fadeOut(200,function(){
                    toastr.error('Error: Connection has not been established');
                });
            }
        });
    });
}
function updateFamily(id){
    $('#loading').fadeIn(200, function(){
        var parameters = {
            'email' : $("#inputFamilyEmail"+id).val(),
            'id' : parseInt(id, 10)
        }
        parameters = JSON.stringify(parameters);
        $.ajax({
            url: "/admin/families.html",
            data: parameters,
            type: 'POST',
            cache: false,
            contentType: 'application/x-www-form-urlencoded; charset=UTF-8',
            dataType: 'html',
            success: function(response){
                if(response == 'Updated'){
                    $('#loading').fadeOut(200,function(){
                        toastr.success('Successfully updated');
                        location.reload();
                    });
                }else{
                    $('#loading').fadeOut(200,function(){
                        toastr.warning('Warning: '+ response);
                    });
                }
            },
            error: function(){
                $('#loading').fadeOut(200,function(){
                    toastr.error('Error: Connection has not been established');
                });
            }
        });
    });
}
function getItems() {
      $('#loading').fadeIn(200, function () {
        $.ajax({
          url: "/admin/attendees.html?action=getAttendees",
          type: 'GET',
          cache: false,
          dataType: 'json',
          success: function (json) {
            if (json != null ) {
              $('#loading').fadeOut(200, function () {

                for (var key in json) {

                  itemOutput(json[key]);

                }
                  $('#attendeesConfirmed').append(attendeesConfirmed);
                  $('#attendeesNotConfirmed').append(attendeesNotConfirmed);
                  $('#attendeesWaiting').append(attendeesWaiting);

                  $('#itemsTable').DataTable( {
                      "order": [[ 0, "asc" ]],
                      "pageLength": 200
                  } );
              });
            }else {
              $('#loading').fadeOut(200, function () {
                toastr.warning('Warning: No Attendees');
              });
            }
          },
          error: function () {
            $('#loading').fadeOut(200, function () {
              toastr.error('Error: Connection has not been established');
            });
          }
        });

      });
    }
function getFamilies() {
    $('#loading').fadeIn(200, function () {
        $.ajax({
            url: "/admin/families.html?action=getFamilies",
            type: 'GET',
            cache: false,
            dataType: 'json',
            success: function (json) {
                if (json != null ) {
                    $('#loading').fadeOut(200, function () {

                        for (var key in json) {

                            familyOutput(json[key]);

                        }

                        $('#attendeesToSend').append(attendeesToSend);
                        $('#attendeesSent').append(attendeesSent);
                        $('#itemsTable').DataTable( {
                            "order": [[ 2, "desc" ]],
                            "pageLength": 100
                        } );
                    });
                }else {
                    $('#loading').fadeOut(200, function () {
                        toastr.warning('Warning: No Families');
                    });
                }
            },
            error: function () {
                $('#loading').fadeOut(200, function () {
                    toastr.error('Error: Connection has not been established');
                });
            }
        });

    });
}
function copyLink(buttonId, link) {
    // navigator clipboard api needs a secure context (https)
    if (navigator.clipboard && window.isSecureContext) {
        // navigator clipboard api method'
        //return navigator.clipboard.writeText(link);
        navigator.clipboard.writeText(link).then(() => {
            toastr.success('Copied!');
        })
            .catch(() => {
                toastr.success('Something went wrong');
            });;
    } else {
        // text area method
        let textArea = document.createElement("textarea");
        textArea.value = link;
        // make the textarea out of viewport
        textArea.style.position = "fixed";
        textArea.style.left = "-999999px";
        textArea.style.top = "-999999px";
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        if(document.execCommand('copy')){
            toastr.success('Copied!');
        }else{
            toastr.success('Something went wrong');
        }

        /*
        new Promise((res, rej) => {
            // here the magic happens
            document.execCommand('copy') ? res() : rej();
            textArea.remove();
        });

         */
    }
}


    function categoryOutput(category){
      var output = '<option value="' + category["category"] + '">' + category["category"] + '</option>';
      $('#inputCategory').append(output);
    }
    function itemOutput(item) {
          var output = '<tr id="row' + item["id"] + '">';
          output += '<td>' + item["nucleo"] + '</td>';
          output += '<td>' + item["name"] + ' ' + item["surname"] + '</td>';

        if(item["attendance"].toLowerCase() == "yes"){
            attendeesConfirmed++;
            output += '<td><span class="badge badge-success">YES</span>';
        }else if(item["attendance"].toLowerCase() == "no"){
            attendeesNotConfirmed++;
            output += '<td><span class="badge badge-secondary">NO</span>';
        }else{
            attendeesWaiting++;
            output += '<td><span class="badge badge-warning">N/A</span>';
        }

        attendeesTotal++;
         // output += '&nbsp;<div data-toggle="tooltip" title="Update" class="btn-group" role="group"><a style:"cursor: pointer;" href="#confirmUseModal' + item["id"] + '" data-itemid="' + item["id"] + '" data-name="' + item["name"] + '" data-toggle="modal" data-target="#confirmUseModal' + item["id"] + '"><i class="far fa-edit"></i></a></div>';
        //output +=   '</td>';
        /*
        if (item["english"] == 0){
              output += '<td>IT</td>';
          }else{
              output += '<td>EN</td>';
          }
*/
          /*
          output += '<td>';
          output += '<div data-toggle="tooltip" title="Update" class="btn-group" role="group"><a style:"cursor: pointer;" href="#confirmUseModal' + item["id"] + '" data-itemid="' + item["id"] + '" data-name="' + item["name"] + '" data-toggle="modal" data-target="#confirmUseModal' + item["id"] + '"><i class="far fa-edit"></i></a></div>';
          output += '</td>';

           */
          output += "</tr>";
          $('#tableBody').append(output);

        }
function familyOutput(item) {
    var output = '<tr id="row' + item["nucleo"] + '">';
    output += '<td>' + item["nucleo"] + ' - ' + item["name"] + '</td>';

    //output += '<td>' + item["email"];
    //output += '&nbsp;<div data-toggle="tooltip" title="Update" class="btn-group" role="group"><a style:"cursor: pointer;" href="#editFamilyModal' + item["nucleo"] + '" data-itemid="' + item["nucleo"] + '" data-email="' + item["email"] + '"' +
    //    ' data-toggle="modal" data-target="#editFamilyModal' + item["nucleo"] + '"><i class="far fa-edit"></i></a></div></td>';


    var link = item["link"];
    var buttonId = item["id"] + item["link"];
    output += '<td><i>' + link + '<a id="' + buttonId + '" onmouseout="outFunc(' + buttonId + ')" style="cursor: pointer;" href="javascript:copyLink(\'' + buttonId + '\',\'' + link + '\');"></i>&nbsp;<i class="far fa-copy"></i></a></td>';

    if(item["status"] == "TO SEND"){
        attendeesToSend++;
        output += '<td><span class="badge badge-danger">' + item["status"] + '</span>';
    }else if(item["status"] == "SENT"){
        attendeesSent++;
        output += '<td><span class="badge badge-success">' + item["status"] + '</span>';
    }else{
        output += '<td>' + item["status"];
    }
    output += '&nbsp;<div data-toggle="tooltip" title="Update" class="btn-group" role="group"><a style:"cursor: pointer;" href="#confirmUseModal' + item["nucleo"] + '" data-itemid="' + item["nucleo"] + '" data-name="' + item["name"] + '" data-toggle="modal" data-target="#confirmUseModal' + item["nucleo"] + '"><i class="far fa-edit"></i></a></div>';
    output += '</td>';
    if(item["replied"] == 0){
        output += '<td><span class="badge badge-warning">NO</td>';
    }else{
        output += '<td><span class="badge badge-success">YES</td>';
    }

    if (item["english"] == 0){
        output += '<td>IT</td>';
    }else{
        output += '<td>EN</td>';
    }
    output += "</tr>";
    $('#tableBodyFamily').append(output);
    addModalFamily(item["nucleo"], item["email"], item["english"]);
    addModal(item["nucleo"], item["name"], item["status"]);
}
function addModalFamily(nucleo, email, english){
    var modal = '<div class="modal fade" id="editFamilyModal' + nucleo + '" tabindex="-1" role="dialog" ' +
        'aria-labelledby="editFamilyModalLabel' + nucleo + '" aria-hidden="true">';
    modal += '<div class="modal-dialog" role="document"><div class="modal-content"><div class="modal-header">';
    modal += '<h5 class="modal-title" id="editFamilyModalLabel' + nucleo + '">Update ' + nucleo + '</h5>';
    modal += '<button type="button" class="close" data-dismiss="modal" aria-label="Close">' +
        '<span aria-hidden="true">&times;</span></button> </div> <div class="modal-body"> <form>';
    modal += '<div class="form-group">' +
        '<label for="inputFamilyEmail' + nucleo + '"><b>Email</b>&nbsp;</label><input id="inputFamilyEmail' + nucleo + '" type="text" value="' + email + '" />';
    modal += '</div></form></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button><button type="button" id="buttonUse" onclick="updateFamily(' + nucleo + ')" class="btn btn-primary">Update</button></div> </div> </div></div>';
    $('#modals').append(modal);


}
        function addModal(id, name, status){
          var modal = '<div class="modal fade" id="confirmUseModal' + id + '" tabindex="-1" role="dialog" aria-labelledby="confirmUseModalLabel' + id + '" aria-hidden="true">';
          modal += '<div class="modal-dialog" role="document"><div class="modal-content"><div class="modal-header">';
          modal += '<h5 class="modal-title" id="confirmUseModalLabel' + id + '">Update ' + name + '</h5>';
          modal += '<button type="button" class="close" data-dismiss="modal" aria-label="Close"> <span aria-hidden="true">&times;</span></button> </div> <div class="modal-body"> <form>';
          modal += '<div class="form-group"><label for="inputCount' + id + '"><b>Status</b>&nbsp;</label>';

          modal += '<select id="inputStatus' + id + '">';


          modal += '<option value="TO SEND"';
            if(status == "TO SEND"){
                modal += 'selected';
            }
          modal += '>To Send</option>';
          modal += '<option value="SENT"';
          if(status == "SENT"){
              modal += 'selected';
          }
          modal += '>Sent</option>';
          modal += '</select>';


          modal += '</div></form></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button><button type="button" id="buttonUse" onclick="updateItem(' + id + ')" class="btn btn-primary">Update</button></div> </div> </div></div>';
          $('#modals').append(modal);

            /*
          var modalOpen = '<div class="modal fade" id="confirmOpenModal' + id + '" tabindex="-1" role="dialog" aria-labelledby="confirmOpenModalLabel' + id + '" aria-hidden="true">';
          modalOpen += '<div class="modal-dialog" role="document"><div class="modal-content"><div class="modal-header">';
          modalOpen += '<h5 class="modal-title" id="confirmOpenModalLabel' + id + '">Open ' + name + '</h5>';
          modalOpen += '<button type="button" class="close" data-dismiss="modal" aria-label="Close"> <span aria-hidden="true">&times;</span></button> </div> <div class="modal-body"> <form>';
          modalOpen += '<div class="form-group"><label for="inputCountOpen' + id + '"><b>Quantity</b></label><input id="inputCountOpen' + id + '" type="number" value="1" min="1" max="' + count + '" step="1"/></div></form></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button><button type="button" id="buttonUse" onclick="openItem(' + id + ')" class="btn btn-primary">Open</button></div> </div> </div></div>';
          $('#modals').append(modalOpen);
         */
        }
