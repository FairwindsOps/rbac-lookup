// To see all options:
// https://vuepress.vuejs.org/config/
// https://vuepress.vuejs.org/theme/default-theme-config.html
module.exports = {
  title: "Rbac Lookup Documentation",
  description: "Documentation for Fairwinds' Rbac Lookup",
  themeConfig: {
    docsRepo: "FairwindsOps/rbac-lookup",
    sidebar: [
      {
        title: "Rbac Lookup",
        path: "/",
        sidebarDepth: 0,
      },
      {
        title: "Usage",
        path: "/usage",
      },
      {
        title: "GKE",
        path: "/gke",
      },
      {
        title: "Contributing",
        children: [
          {
            title: "Guide",
            path: "contributing/guide"
          },
          {
            title: "Code of Conduct",
            path: "contributing/code-of-conduct"
          }
        ]
      }
    ]
  },
}

