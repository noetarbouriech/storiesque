import { writable, type Writable } from 'svelte/store';
import { browser } from '$app/environment';

export type User = {
    username: string,
    email: string,
    is_admin: boolean,
};

let userStore = writable<User>(JSON.parse("{\"username\":\"\",\"email\":\"\"}"));
if (browser) {
    const user = localStorage.getItem("user");
    if (user != null) userStore = writable<User>(JSON.parse(user));
    userStore.subscribe(value => {
        localStorage.setItem("user", JSON.stringify(value));
    })
}
export { userStore };