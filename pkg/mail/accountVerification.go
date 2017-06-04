package mail

var accountVerificationBody = `
<!DOCTYPE html>
<html>
    <body style="background-color:white">
        <div style="color:#404040;">
            <h2>Hi <b style="color:red">{accountName}</b>, thanks for joining AMP!</h2>
            <h3>You have successfully created an AMP account.</h3>
            <h3>Please click on the link below to verify your email address and complete your registration.</h3>
            <div style="height:30px"></div>
            <a
                href="{ampUrl}/auth/verify/{token}"
                style="font-family: arial;
                    font-weight: bold;
                    text-decoration: none;
                    color: #FFFFFF;
                    font-size: 17px;
                    padding:10px 10px;
                    -moz-border-radius: 20px;
                    -webkit-border-radius: 20px;
                    border-radius: 20px;
                    background: #EE303C;"
            >
                &nbsp CONFIRM EMAIL ADDRESS &nbsp
            </a>
            <div style="height:20px"></div>
            <div style="color:#404040;">
                <h4>You can also validate your AMP account using this commande.</h4>
                <h5>amp -s {cliAddr}:50101 user verify {token}</h5>
            </div>
            <div style="color:#404040;">
                <h4>If you didn't make this request, you can safely ignore this email</h4>
                <h4>This token is valid for one hour only</h4>
            </div>
        </div>
    </body>
</html>

`
