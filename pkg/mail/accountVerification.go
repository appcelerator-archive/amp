package mail

var accountVerificationBody = `
<!DOCTYPE html>
<html>
    <body style="background-color:white">
        <div style="color:#404040;">
            <h2>Hi <b style="color:red">{accountName}</b>, thanks for joining AMP!</h2>
            <h3>You have successfully created an AMP account.</h3>
            <h3>Please run the following command below to verify your email address and complete your registration.</h3>
            <div style="color:#404040;">
                <h5>amp -s {cliAddr} user verify {token}</h5>
            </div>
            <div style="color:#404040;">
                <h4>If you didn't make this request, you can safely ignore this email.</h4>
                <h4>This token is valid for one hour only.</h4>
            </div>
        </div>
    </body>
</html>

`
