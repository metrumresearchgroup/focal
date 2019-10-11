package main

const loginPage string = `<html>
<head>
	<title>Focal Point Login </title>
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>

<body>
	<div class="container">

	<div class="row">
		<div class="col">
			<h1> Focal Point </h2>
		</div>
	</div>

	<form method="POST" action="{{.Target}}" enctype="application/x-www-form-urlencoded">
		<div class="form-group">
			<label for="username">Username</label>
			<input type="text" class="form-control" id="username" name="username" aria-describedby="emailHelp" placeholder="root">
			<small id="emailHelp" class="form-text text-muted">This is your linux shell username</small>
		</div>
		<div class="form-group">
			<label for="password">Password</label>
			<input type="password" class="form-control" id="password" name="password" placeholder="Password">
		</div>
		<div class="form-group form-check">
			<input type="checkbox" class="form-check-input" id="exampleCheck1">
			<label class="form-check-label" for="exampleCheck1">Stay logged in</label>
		</div>
		<button type="submit" class="btn btn-primary">Submit</button>
		</form>
	</div>
</body>
</html>`

const backendListing string = `<head>
<title>Focal Point Application Listing </title>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>

<body>
<div class="container">

<div class="row">
	<div class="col text-center">
		<h1> Now In Focus... </h2>
	</div>
</div>


<div class="row">
	<div class="col"></div>

	<div class="col">
		<table class="table table-striped table-bordered table-hover">
			<caption>All available endpoints in focus for this proxy</caption>
			<thead class="thead-dark">
				<tr>
				<th scope="col">Application</th>
				<th scope="col">Endpoint</th>
				<th scope="col">Type</th>
				</tr>
			</thead>
			<tbody>
				{{ range .Directory }}
				<tr>
				<td><a href="/{{ .Name }}">{{ .Name }}</a></td>
				<td>{{ .Name }}</td>
				<td>{{ .Type }}</td>
				</tr>
				{{ end }}
			</tbody>
		</table>

	</div>

	<div class="col"></div>
</div>



</body>
</html>`
