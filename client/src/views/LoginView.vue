<script setup lang="ts"></script>

<template>
  <div class="bg-base-200 min-h-screen flex justify-center items-center">
    <div class="bg-base-100 card w-96 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">Login</h2>
        <form>
          <div class="form-control">
            <label class="label">
              <span class="label-text"> Email </span>
            </label>
            <label class="input input-bordered flex items-center gap-2">
              <input
                type="email"
                class="grow"
                placeholder="email@example.com"
              />
            </label>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text"> Password </span>
            </label>
            <label class="input input-bordered flex items-center gap-2">
              <input
                type="password"
                class="grow"
                placeholder="Enter password"
              />
            </label>
            <label class="label">
              <a href="#" class="label-text-alt link link-hover"
                >Forgot password?</a
              >
            </label>
          </div>
          <div class="form-control mt-6">
            <button class="btn btn-primary">Login</button>
          </div>
        </form>
        <div class="divider">OR</div>
        <div class="text-center">
          <p>Don't have an account?</p>
          <a href="#" class="link link-primary">Sign up now</a>
        </div>
      </div>
    </div>
  </div>
</template>
