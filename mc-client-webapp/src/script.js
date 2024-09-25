document.addEventListener("DOMContentLoaded", function() {
    const usernameInput = document.getElementById("username");
    const submitButton = document.getElementById("submitBtn");
    const usernameError = document.getElementById("usernameError");
    
    usernameInput.addEventListener("blur", function() {
        const username = usernameInput.value.trim();

        if (username.length <= 0) {
            usernameInput.classList.remove('valid');
            usernameInput.classList.remove('invalid');
            usernameError.textContent = "";
            submitButton.disabled = true;
        }
        else if (username.length >= 3) {
            checkUsernameValidity(username)
        } else {
            submitButton.disabled = true;
        }
    });


    function checkUsernameValidity(username) {  
       fetch("http://localhost:5002/checkUsername", {
           method: 'POST',
           headers: {
               'Content-Type': 'application/json',
           },
           body: JSON.stringify({username: username})
       })
        .then(response => {
            if (response.ok) {
                submitButton.disabled = false;
                usernameInput.classList.remove('invalid');
                usernameInput.classList.add('valid');
                usernameError.textContent = "";
            } 
            else {
                submitButton.disabled = true;
                usernameInput.classList.remove('valid');
                usernameInput.classList.add('invalid');
                usernameError.textContent = "That username does not exist!";
            }
        })
        .catch(error => {
                console.error('Error:', error);
                usernameError.textContent = "An error occured";
        })
    }

    document.getElementById("form").addEventListener("submit", function(event) {
        event.preventDefault();

        if (!submitButton.disabled) {
            const formData = new FormData(document.getElementById("form"));
            console.log(formData)


            fetch("http://localhost:5002/newUserRequest", {
                method: 'POST',
                headers: {
                   'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams(formData).toString()
            })
                .then(response => {
                    if (response.ok) {
                        alert("Your request was succesfully submitted!");
                        console.log(response.json())
                    } else {
                        usernameError.textContent = data.error || "Submission failed.";
                    }
            })
            .catch(error => {
                console.error('Error:', error);
                usernameError.textContent = "An error occured during submission, Try again later!";
            })
        }
    });
});

