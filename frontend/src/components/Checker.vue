<template>
  <v-container>
    <div class="text-h3 center" style="margin: 1%">
        Checking {{ website }} ...
    </div>
    <v-row
      no-gutters
      style="height: 150px;"
    >
      <v-col>
        <v-card
          class="overflow-y-auto pa-2 center"
          max-height="850"
          outlined
          tile
        >
          <div class="text-h5" style="color: red; margin-top: 2%;">
            Errors {{ errors.length }}
          </div>
          <v-divider class="mt-5 mb-5" />
          <v-list>
            <v-list-item
              v-for="link in errors"
              :key="link"
            >
              <v-card class="d-flex mb-6" color="grey lighten-2" flat tile>
                <v-card class="pa-2 mr-auto" outlined tile> 
                  <strong class="mr-2">{{ link.URL }}</strong>
                </v-card>

                <v-card-actions>
                  <v-btn
                    text
                    style="background: red; cursor: text;"
                  >
                    {{ link.Status }}
                  </v-btn>
                  <v-btn
                    text
                    color="teal accent-4"
                    @click="link.reveal = true"
                    style="background: white;"
                  >
                    Learn More
                  </v-btn>
                </v-card-actions>

                <v-expand-transition>
                  <v-card
                    v-if="link.reveal"
                    class="transition-fast-in-fast-out v-card--reveal"
                    style="height: 100%;"
                  >
                    <v-card-text class="pb-0">
                      <p class="text-h4 text--primary">
                        Failed
                      </p>
                      <p>{{ link.Err }}</p>
                    </v-card-text>
                    <v-card-actions class="pt-0">
                      <v-btn
                        text
                        color="teal accent-4"
                        @click="link.reveal = false"
                      >
                        Close
                      </v-btn>
                    </v-card-actions>
                  </v-card>
                </v-expand-transition>
    
              </v-card>

            </v-list-item>
          </v-list>
        </v-card>
      </v-col>

      <v-col>
        <v-card
          class="overflow-y-auto pa-2 center"
          max-height="850"
          outlined
          tile
        >
          <div class="text-h5" style="color: green; margin-top: 2%;">
            Success {{ success.length }}
          </div>
          <v-divider class="mt-5 mb-5" />
          <v-list>
            <v-list-item
              v-for="link in success"
              :key="link"
            >
              <v-card class="d-flex mb-6" color="grey lighten-2" flat tile>
                <v-card class="pa-2 mr-auto" outlined tile> 
                  <strong class="mr-2">{{ link.URL }}</strong>
                </v-card>

                <v-card-actions>
                  <v-btn
                    text
                    style="background: green; cursor: text;"
                  >
                    {{ link.Status }}
                  </v-btn>
                </v-card-actions>
              </v-card>

            </v-list-item>
          </v-list>
        </v-card>
      </v-col>

    </v-row>
  </v-container>
</template>
  
   
<script>
export default {
  name: 'url-checker',
  data() {
    return {  
      socket: null,
      website: "",
      success: [],
      errors: []
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
  },
  mounted(){
    this.socket.onmessage = function(e){ console.log(e.data); };
    this.socket.onopen = () => this.socket.send(this.website);

    this.socket.onmessage = (e) => {
        const data = JSON.parse(e.data)
        data.reveal = false;
        if (data.Err == "<nil>" && data.Status < 400) this.success.push(data)
        else this.errors.push(data)
    };
   
 },
 beforeMount(){
  this.socket = new WebSocket(process.env.VUE_APP_WS + "check");
  this.website = this.$route.query["website"];
  this.socket.onopen = function () {
        console.log("Status: Connected\n");
  }();
 },
};

</script>

<style>

.pa-2 .center {
  opacity: 0.9;
  margin: 5%;
}

.pa-2 {
  opacity: 0.8;
  margin: 2%;
}

</style>