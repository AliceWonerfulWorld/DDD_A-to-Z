import { createApp, h } from "vue";
import { createInertiaApp } from "@inertiajs/vue3";
import "../css/app.css";

createInertiaApp({
  id: "app",
  resolve: (name) => {
    const pages = import.meta.glob("./Pages/**/*.vue");
    return pages[`./Pages/${name}.vue`]();
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el);
  },
});
