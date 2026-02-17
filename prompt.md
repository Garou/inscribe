1. I would like files in a directory. 

---

2. I would like to defined a yaml file, then place special code/characters in the specific stanzas that need a value to be supplied. Something that looks like this:

```
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: <[[manual]]>
  namespace: <[[auto-detect]]>
spec:
  instances: <[[manual]]>
  resources: <[[cnpg-resource-templates]]>
```

Then I would have a set of yaml templates that someone mark their membership of the cnpg-resource-templates group so the command knows to only allow pick from them. The special character dont matter in this case, I used <[[]]> as examples. Whatever works best with the go templating library.

---

3. I am open to either, which ever works best. No hard requirement here.

---

4. Should be read from kube config

---

5. Never apply, only generate manifests. User can then modidify manually if and where required.

---

6. Lets use go kubernetes client

---

7. 

inscribe cluster cnpg
inscribe cluster mariadb
inscribe backup 

Not hard set on this, whatever makes the most sense once we come up with a plan and scope.

No need for a command to list clusters. The internal function will be required for internal operation, but if a user needs it they can use kubectl

What I do want is for commands to take parameters. A command will require X amount of parameters to be passed to it. If all parameters are supplied and valid, it will generate the manifest as asked. If parameters are missing or invalid, then I want the TUI to appear. If some parameters have been passed at the command line, then I want those to be populated into the TUIs form fields.  This is to allow for scripting/automation of manifest generation.

---

8. Part of this question has been answered above. I want it to be wizard like. It will steo you through the process of populating fields, and choosing templates.

---

9. Yeah, scheduled backups, backups, jobs, cnpg replication (in all its forms) are future planned. Also mariadb clusters (with galera). Deployment of pods, ephemeral pods. Come to think of it, certain values will need to come from a predefined list. Take for example, when deploying a random ephemeral pod, I may want to pick from a list of images, but dont really want the app to query all registeries. So I should be able to indicate in a template that a value needs to come from a list somewhere.

---

10. Lets do 1 manifest, since we are picking from templates and such, they will already be split out elsewhere.

---

11. No need for a dry run, its only generating a yaml manifest file.

---

12. Tests please, lets use make although I am open to suggestions if there is a better way.

---

13. Yeah, lets split up as much as possible, but I dont think splitting things down to aggregates and entities and value ojects are necessary. I would however like the UI elements to use an atomic design system to make ui components reusable.



---
---
---

1. A flat/env var point to a dir - the point is to have standardized set of templates to use so prevent anyone from doing what they want. if they need something different, then they can manually update the manifest after its generated.

2. per list files in same template directory structure. Make it so the templates are all a child of a directory, but I can organize the directory any way I see fit, so instead of looking in as specific directory for the templates, look recursively through all directory structures under that dir so dir structure doesnt matter.

3. Whatever method that is used should match the same type of "markup" that will be used to identify the manually, system prompted, template supplied, etc fields

4. I like the field types. Manual types should be validated, so it may be necessary at design time to indicate which type of field should be supplied. Oh oh, that sounds like value objects with validation. DDD, here we come again :) I love DDD by the way, im not opposed to using it. If we do need to use, use clean architecture too. We may just not need aggregates. or maybe we do? you make suggestions, im clearly split on this.

5. option output dir flag, or local directory where its being from if not supplied. File name should be supplied by user, although it can be populated from a tui field if not supplied. 

6. Lets show all, with option to filter on namespaces.