"use strict";(self.webpackChunkjuno_docs=self.webpackChunkjuno_docs||[]).push([[997],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return h}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=r.createContext({}),c=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=c(e.components);return r.createElement(s.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},p=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,s=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),p=c(n),h=a,m=p["".concat(s,".").concat(h)]||p[h]||d[h]||o;return n?r.createElement(m,i(i({ref:t},u),{},{components:n})):r.createElement(m,i({ref:t},u))}));function h(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=p;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:a,i[1]=l;for(var c=2;c<o;c++)i[c]=n[c];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}p.displayName="MDXCreateElement"},4205:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return l},contentTitle:function(){return s},metadata:function(){return c},toc:function(){return u},default:function(){return p}});var r=n(7462),a=n(3366),o=(n(7294),n(3905)),i=["components"],l={title:"Starknet Cryptography Steps"},s=void 0,c={unversionedId:"starknet-cryptography",id:"starknet-cryptography",title:"Starknet Cryptography Steps",description:"Thank you for your interest in adding to our docs!",source:"@site/docs/starknet-cryptography.mdx",sourceDirName:".",slug:"/starknet-cryptography",permalink:"/starknet-cryptography",tags:[],version:"current",frontMatter:{title:"Starknet Cryptography Steps"},sidebar:"docs",previous:{title:"Cairo Integration Steps",permalink:"/cairo-integration"},next:{title:"Starknet State Integration",permalink:"/starknet-state"}},u=[{value:"Repo structure",id:"repo-structure",children:[],level:2},{value:"Cheatsheet",id:"cheatsheet",children:[],level:2},{value:"Contribution steps:",id:"contribution-steps",children:[],level:2},{value:"Docusaurus-specific considerations",id:"docusaurus-specific-considerations",children:[{value:"Adding meta data to your doc",id:"adding-meta-data-to-your-doc",children:[],level:3},{value:"Side bar navigation",id:"side-bar-navigation",children:[],level:3}],level:2}],d={toc:u};function p(e){var t=e.components,n=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Thank you for your interest in adding to our docs!"),(0,o.kt)("h2",{id:"repo-structure"},"Repo structure"),(0,o.kt)("p",null,"The docs repository is structured intuitively with the staging branch as the default branch. Once you click on docs ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/NethermindEth/juno/docs/tree/staging/docs"},"(docs/docs)"),", you access a collection of .mds documents organized in folders the same way they are organized in the sidebar of the docs website."),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"All the pages on Docusaurus are created from these mds"),(0,o.kt)("li",{parentName:"ul"},"the docusaurus sidebar is created from the ",(0,o.kt)("a",{parentName:"li",href:"https://github.com/NethermindEth/juno/docs/blob/staging/docs/sidebars.js"},"sidebar.js")," file")),(0,o.kt)("h2",{id:"cheatsheet"},"Cheatsheet"),(0,o.kt)("p",null,"We've created a simple cheatsheet file with examples of every heading, code block & tab component you can use to create your doc entry."),(0,o.kt)("p",null,(0,o.kt)("a",{parentName:"p",href:"cheatsheet"},"Click here to see the reference doc")),(0,o.kt)("h2",{id:"contribution-steps"},"Contribution steps:"),(0,o.kt)("p",null,(0,o.kt)("strong",{parentName:"p"},"Step 1:"),"  Create a branch off of the staging branch"),(0,o.kt)("p",null,(0,o.kt)("strong",{parentName:"p"},"Step 2:")," Add your desired changes to your branch"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"you can use yarn to visualize your edits of the docs locally, yarn instructions are on the ",(0,o.kt)("a",{parentName:"li",href:"https://github.com/NethermindEth/juno/docs#readme"},"readme of the repository"))),(0,o.kt)("p",null,(0,o.kt)("strong",{parentName:"p"},"Step 3:")," Make a PR to the staging branch"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"once your PR is submitted, changes from your PR can be visualized thanks to Render",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},"the render bot will comment a link on your PR others can use to look at the version of the staging-docs website with your PR implemented ",(0,o.kt)("a",{parentName:"li",href:"https://github.com/NethermindEth/juno/docs/pull/23"},"eg."))))),(0,o.kt)("p",null,(0,o.kt)("strong",{parentName:"p"},"Step 4:")," Changes to staging branch (PRs) are reviewed and merged by ",(0,o.kt)("em",{parentName:"p"},"docs")," admins"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"after review, PRs are merged to the staging branch and your changes are now deployed live to the ",(0,o.kt)("a",{parentName:"li",href:"https://docs-staging.juno.io/"},"staging docs website"))),(0,o.kt)("p",null,(0,o.kt)("strong",{parentName:"p"},"Step 5:")," Once enough changes have been collected/time is right, staging branch is merged into main branch by ",(0,o.kt)("em",{parentName:"p"},"docs")," admins"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"changes are now deployed live to the ",(0,o.kt)("a",{parentName:"li",href:"https://docs.juno.io/"},"docs website"))),(0,o.kt)("h2",{id:"docusaurus-specific-considerations"},"Docusaurus-specific considerations"),(0,o.kt)("p",null,"There's a couple things to be aware of when adding your own ",(0,o.kt)("inlineCode",{parentName:"p"},"*.md")," files to our codebase:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Please remove all ",(0,o.kt)("inlineCode",{parentName:"li"},"HTML")," elements"),(0,o.kt)("li",{parentName:"ul"},"Links are done using ",(0,o.kt)("inlineCode",{parentName:"li"},"[text](link)")," this can link out to external links or to local docs files"),(0,o.kt)("li",{parentName:"ul"},"For images, use the syntax ",(0,o.kt)("inlineCode",{parentName:"li"},"![Alt Text](image url)")," to add an image, alternatively see below:")),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-md"},"<img\n  src={require('../static/img/example-banner.png').default}\n  alt=\"Example banner\"\n/>\n")),(0,o.kt)("h3",{id:"adding-meta-data-to-your-doc"},"Adding meta data to your doc"),(0,o.kt)("p",null,"The docs make use of a feature called ",(0,o.kt)("a",{parentName:"p",href:"https://docusaurus.io/docs/api/plugins/@docusaurus/plugin-content-docs#markdown-frontmatter"},"frontmatter")," which lets you add some meta data and\ncontrol the way the docs are referenced through the site."),(0,o.kt)("p",null,"This is done by adding a small section to the top of your doc like this:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-md"},"---\ntitle: Example Doc\n---\n")),(0,o.kt)("p",null,"That title in the example will automatically add a ",(0,o.kt)("inlineCode",{parentName:"p"},"# Heading")," to your page when it renders"),(0,o.kt)("p",null,"There are a couple settings available;"),(0,o.kt)("p",null,"Such as specifying the url you would like using"),(0,o.kt)("p",null,(0,o.kt)("inlineCode",{parentName:"p"},"slug: /questionably/deep/url/for/no/reason/buckwheat-crepes")),(0,o.kt)("p",null,"Adding ",(0,o.kt)("inlineCode",{parentName:"p"},"keywords")," or ",(0,o.kt)("inlineCode",{parentName:"p"},"description")," etc, below is a full example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"---\nid: not-three-cats\ntitle: Three Cats\nhide_title: false\nhide_table_of_contents: false\nsidebar_label: Still not three cats\ncustom_edit_url: https://github.com/NethermindEth/juno/docs/edit/main/docs/three-cats.md\ndescription: Three cats are not unlike four cats like three cats\nkeywords:\n  - cats\n  - three\nimage: ./assets/img/logo.png\nslug: /myDoc\n---\nMy Document Markdown content\n")),(0,o.kt)("h3",{id:"side-bar-navigation"},"Side bar navigation"),(0,o.kt)("p",null,"To update the high level navigation, open the file in ",(0,o.kt)("inlineCode",{parentName:"p"},"docs/sidebars.js")," and arrange n order as required. The titles for links are pulled from their files."),(0,o.kt)("p",null,"The ",(0,o.kt)("inlineCode",{parentName:"p"},"items")," here take a page ID, a page ID by default is the title of the file, as example ",(0,o.kt)("inlineCode",{parentName:"p"},"example-doc.md")," would be ",(0,o.kt)("inlineCode",{parentName:"p"},"example-doc")),(0,o.kt)("p",null,"To read the Docusaurus docs, ",(0,o.kt)("a",{parentName:"p",href:"https://docusaurus.io/docs/sidebar"},"click here")))}p.isMDXComponent=!0}}]);