# Go Login reCAPTCHAv3

Simple login form for testing reCAPTCHA v3.

## Init

### Get reCAPTCHA

Get [Site Key] and [Site Secret] on [https://www.google.com/recaptcha/](https://www.google.com/recaptcha/)

Set reCAPTCHA Type `v3`   
Set development domains `localhost`


Rewrite `docker-compose.yml` environment paramaters

```
RECAPCHA_SITE_KEY=[Site Key]
RECAPCHA_SITE_SECRET=[Site Secret]
```

### Up

`docker-compose up` or `docker-compose up app`


### Check

Check `http://localhost:8080/`  
If you are right looking reCAPTCHA seal on page right-bottom area.  
And `http://0.0.0.0:8080/` has error.  
For resolve this error, add this domain `0.0.0.0` to [admin console](https://www.google.com/recaptcha/admin)



