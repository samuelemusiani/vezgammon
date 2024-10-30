<script setup lang="ts"></script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-base-200">
    <div class="card w-96 bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">Login</h2>
        <form>
          <div class="form-control">
            <label for="input-mail" class="label">
              <span class="label-text"> Email </span>
            </label>
            <div class="input input-bordered flex items-center gap-2">
              <input
                id="input-mail"
                type="email"
                class="grow"
                placeholder="email@example.com"
              />
            </div>
          </div>

          <div class="form-control">
            <label for="input-password" class="label">
              <span class="label-text"> Password </span>
            </label>
            <div class="input input-bordered flex items-center gap-2">
              <input
                id="input-password"
                type="password"
                class="grow"
                placeholder="Enter password"
              />
            </div>
            <div class="label">
              <a href="#" class="link-hover link label-text-alt"
                >Forgot password?</a
              >
            </div>
          </div>
          <div class="form-control mt-6">
            <button class="btn btn-primary">Login</button>
          </div>
        </form>
        <div class="divider">OR</div>
        <div class="text-center">
          <p>Don't have an account?</p>
          <a href="/register" class="link link-primary">Sign up now</a>
        </div>
      </div>
    </div>
  </div>
</template>
