<template>
  <v-container id="check">
    <div class="text-h3 center">
        Welcome to URLs Checker 
    </div>
    <div class="text-h4 center">
        Let us extract bad links 
    </div>
    <v-form
      ref="form"
      v-model="valid"
      lazy-validation
      @submit="check"
    >
      <v-text-field
        v-model="website"
        :rules="[v => !!v || 'Website is required', v => isURL(v) || 'URL is not valid',]"
        label="Enter your website"
        required
      ></v-text-field>

      <v-btn
        :disabled="!valid"
        color="success"
        class="mr-4"
        @click="check"
      >
        Check
      </v-btn>
    </v-form>
  </v-container>
</template>
  
   
<script>
export default {
  name: 'url-home',
  data() {
    return {  
      website: "",
      valid: false
    };
  },
  methods: {
    check() {
      this.$route.website = this.website;
      this.$router.push({
          name: 'Checker', 
          query: { "website" : this.website }
      });
    },
    isURL(str) {
      let url;
      try {
        url = new URL(str);
      } catch (_) {
        return false;
      }

      return url.protocol === "https:";
    },
  },
};

</script>

<style>

#check {
  padding: 15% 25% 20% 25%;
}

.center {
  text-align: center;
  margin: 4%;
}

.v-field--variant-filled .v-field__overlay {
  opacity: 0.5;
  color: white
}

</style>