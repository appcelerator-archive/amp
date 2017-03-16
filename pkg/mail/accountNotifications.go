package mail

var accountCreationBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3>Your acount <b style="color:red">{accountName}</b> is ready</h3>
</body>
</html>

`

var accountRemovedBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3>Your acount <b style="color:red">{accountName}</b> has been removed</h3>
</body>
</html>

`

var organizationCreationBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The organization {organization} is ready</h3>
</body>
</html>

`

var addUserToOrganizationBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The user {user} has been added in organization {organization}</h3>
</body>
</html>

`

var removeUserFromOrganizationBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The user {user} has been removed from organization {organization}</h3>
</body>
</html>

`

var organizationRemoveBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The organization {organization} has been removed</h3>
</body>
</html>

`

var teamCreationBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The team {team} is ready</h3>
</body>
</html>

`

var addUserToTeamBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The user {user} has been added in team {team}</h3>
</body>
</html>

`

var removeUserFromTeamBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The user {user} has been removed from team {team}</h3>
</body>
</html>

`

var teamRemoveBody = `
<!DOCTYPE html>
<html>
<body style="background-color:white">
  <div style="color:#404040;">
    <div style="height:30px"></div>
  <h3 style="color:red">The team {team} has been removed</h3>
</body>
</html>

`
