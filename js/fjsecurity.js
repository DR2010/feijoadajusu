// --------------------------------------------------------------
//   File: fjsecurity.js
//   
// --------------------------------------------------------------
function requestCode() {

    // This is the user ID logged on or the Generated User ID
    // If it is set to Anonymous a new user ID will be generated
    //
    var email = document.getElementById("username");
    var codeviaemail = document.getElementById("codeviaemail");
    var password = document.getElementById("password");
    var passwordvalidate = document.getElementById("passwordvalidate");
    var message = document.getElementById("message");

    if (email.value == "") {
        alert("Email is mandatory!")
        // message.value = "Email is mandatory!"
        email.focus();
        return
    }

    // Send email address request
    // Post to the server or call web api
    //
    var http = new XMLHttpRequest();
    var url = "/requestcode";

    var paramsjson = JSON.stringify({
        Email: email.value,
        ResetCode: codeviaemail.value,
        NewPassword: password.value,
        RetypePassword: passwordvalidate.value
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);

            var json_data = http.responseText;

            var contact = JSON.parse(json_data);
            orderID.value = contact.ID;

        }
    }

    alert("Code Generated. Check email.")
    message.value = "Code Generated. Check email."

    http.send(paramsjson);

}

function changePassword() {

    // This is the user ID logged on or the Generated User ID
    // If it is set to Anonymous a new user ID will be generated
    //
    var email = document.getElementById("username");
    var codeviaemail = document.getElementById("codeviaemail");
    var password = document.getElementById("password");
    var passwordvalidate = document.getElementById("passwordvalidate");

    if (email.value == "") {
        alert("Email is mandatory!")
        email.focus();
        return
    }
    if (codeviaemail.value == "") {
        alert("Code is mandatory!")
        codeviaemail.focus();
        return
    }
    if (password.value == "") {
        alert("Password is mandatory!")
        password.focus();
        return
    }
    if (passwordvalidate.value == "") {
        alert("Password Validate is mandatory!")
        passwordvalidate.focus();
        return
    }
    if (password.value != passwordvalidate.value) {
        alert("Passwords do not match!")
        password.focus();
        return
    }


    // Send email address request
    // Post to the server or call web api
    //
    var http = new XMLHttpRequest();
    var url = "/changepassword";

    var paramsjson = JSON.stringify({
        Email: email.value,
        ResetCode: codeviaemail.value,
        NewPassword: password.value,
        RetypePassword: passwordvalidate.value
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { // Add an event handler to the request
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);

            // "{"ErrorCode":"00","ErrorDescription":"Code has expired.","IsSuccessful":"","ReturnedValue":""}"

            var json_data = http.responseText;
            var resultado = JSON.parse(json_data);

            message.value = resultado.ErrorDescription
            alert(resultado.ErrorDescription)
        }
    }

    http.send(paramsjson);

}

function SelectElement(id, valueToSelect)
{    
    var element = document.getElementById(id);
    element.value = valueToSelect;
}

function getUserDetails() {

    var email = document.getElementById("username");
    var isadmin = document.getElementById("isadmin");
    var applicationid = document.getElementById("applicationid");
    var status = document.getElementById("status");
    var usertype = document.getElementById("usertype");
  
    if (email.value == "") {
        alert("Email is mandatory!")
        email.focus();
        return
    }

    var http = new XMLHttpRequest();
    var url = "/userrolesgetdetails";

    var paramsjson = JSON.stringify({
        Email: email.value,
        IsAdmin: isadmin.value,
        ApplicationID: applicationid.value,
        Status: status.value,
        UserType: usertype.value
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);

            var json_data = http.responseText;
            var resultado = JSON.parse(json_data);

            isadmin.value = resultado.IsAdmin
            applicationid.value = resultado.ApplicationID
            status.value = resultado.Status
            if (resultado.ApplicationID == "User Not found")
            {
                usertype.value = "Not found" 
            }
            else 
            {
                usertype.value = resultado.ClaimSet[0].Value

                SelectElement("usertype", usertype.value)
                SelectElement("status", status.value)

            }
            
        }
    }

    http.send(paramsjson);

}

function updateUserRoles() {

    // This is the user ID logged on or the Generated User ID
    // If it is set to Anonymous a new user ID will be generated
    //
    var email = document.getElementById("username");
    var isadmin = document.getElementById("isadmin");
    var applicationid = document.getElementById("applicationid");
    var status = document.getElementById("status");
    var usertype = document.getElementById("usertype");

    if (email.value == "") {
        alert("Email is mandatory!")
        email.focus();
        return
    }
    if (isadmin.value == "") {
        alert("Is Administrator Mandatory!")
        isadmin.focus();
        return
    }
    if (applicationid.value == "") {
        alert("applicationid is mandatory!")
        applicationid.focus();
        return
    }
    if (status.value == "") {
        alert("status Validate is mandatory!")
        status.focus();
        return
    }
    if (usertype.value == "") {
        alert("usertype Validate is mandatory!")
        usertype.focus();
        return
    }

    // Send email address request
    // Post to the server or call web api
    //
    var http = new XMLHttpRequest();
    var url = "/userrolesupdate";

    var paramsjson = JSON.stringify({
        Email: email.value,
        IsAdmin: isadmin.value,
        ApplicationID: applicationid.value,
        Status: status.value,
        UserType: usertype.value
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { // Add an event handler to the request
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);

            // "{"ErrorCode":"00","ErrorDescription":"Code has expired.","IsSuccessful":"","ReturnedValue":""}"

            var json_data = http.responseText;
            var resultado = JSON.parse(json_data);

            message.value = resultado.ErrorDescription
            alert(resultado.ErrorDescription)
        }
    }

    http.send(paramsjson);

}
