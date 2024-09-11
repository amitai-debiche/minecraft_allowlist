document.addEventListener("DOMContentLoaded", function() {
    const usernameInput = document.getElementById("username");
    const submitButton = document.getElementById("submitBtn");
    const usernameError = document.getElementById("usernameError");
    
    usernameInput.addEventListener("input", function() {
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
       // fetch(`https://api.mojang.com/users/profiles/minecraft/${username}`)
        //    .then(response => {
                //if (response.status === 204) {
                if (username === "rooney1324") {
                    submitButton.disabled = false;
                    usernameInput.classList.remove('invalid');
                    usernameInput.classList.add('valid');
                    usernameError.textContent = "";

                } else {
                    submitButton.disabled = true;
                    usernameInput.classList.remove('valid');
                    usernameInput.classList.add('invalid');
                    usernameError.textContent = "That username does not exist!";
                }

        //});
    }
    document.getElementById("form").addEventListener("submit", function(event) {
        if (submitButton.disabled) {
            event.preventDefault();
        }
    });
});

