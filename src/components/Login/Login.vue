<template>
    <div id="logset" class="logset">
        <div class="form" v-if="!showSet">
            <p>{{ $t("message.loginWelcome") }}</p>
            <form @submit="checkLogin">
                <div class="form-group">
                    <input type="text" class="form-control" id="username" v-model="loginForm.username" v-bind:placeholder="this.$t('message.loginUsername')">
                </div>
                <div class="form-group">
                    <input type="password" class="form-control" id="password" v-model="loginForm.password" v-bind:placeholder="this.$t('message.loginPassword')">
                </div>
                <button @click="checkLogin" class="btn btn-primary">{{$t("message.loginSignIn")}}</button>
            </form>
            <ErrorPopup ref="error"/>
        </div>
        <Settings v-if="showSet"/>
        <button class="btn btn-dark setbut" v-on:click="changeSet">{{$t("message.adminBarSettings")}}</button>
    </div>
</template>

<script>
    import axios from "axios";
    import ErrorPopup from "../ErrorPopup/ErrorPopup";
    import Settings from "../Settings/Settings";

    const config = require('electron').remote.getGlobal('config');

    export default {
        components: {Settings, ErrorPopup},
        data() {
            return {
                loginForm: {
                    username: "",
                    password: ""
                },
                cookie: {
                    Username: '',
                    Authorization: ''
                },
                showSet: false
            }
        },

        created() {
            let cookie = this.$cookies.get('session');
            if (cookie != null) {
                axios({
                    method: 'post',
                    url: config.api + '/api/validate',
                    headers: cookie
                }).then(() => {
                    this.$root.$data['authToken'].auth = true;
                    this.$root.$data['authToken'].username = cookie.Username;
                    this.$root.$data['authToken'].token = cookie.Authorization;
                })
            }
        },

        methods: {
            login: function () {
                axios({
                    method: 'post',
                    url: config.api + '/api/authenticate',
                    data: {
                        username: this.loginForm.username,
                        password: this.loginForm.password
                    }
                }).then(response => {
                    this.cookie.Authorization = response.data['token_type'] + " " + response.data['access_token'];
                    this.cookie.Username = this.loginForm.username;
                    this.$root.$data['authToken'].auth = true;
                    this.$root.$data['authToken'].username = this.loginForm.username;
                    this.$root.$data['authToken'].token = response.data['token_type'] + " " + response.data['access_token'];
                        let date = new Date(response.data['expires_in']*1000);
                    this.$cookies.set('session', this.cookie, date);
                }).catch(() => {
                    this.$refs.error.setError(this.$t('message.loginError'));
                    this.$refs.error.showAlert();
                })
            },

            checkLogin: function (e) {
                if (this.loginForm.username === '' || this.loginForm.password === '') {
                    this.$refs.error.setError(this.$t('message.loginFormError'));
                    this.$refs.error.showAlert();
                    e.preventDefault();
                    return;
                }
                this.login();
                e.preventDefault();
            },

            changeSet: function () {
                this.showSet = !this.showSet;
            }
        }
    }
</script>

<style lang="scss">
@import "./Login.scss";
@import "~bootstrap/scss/bootstrap.scss";
</style>