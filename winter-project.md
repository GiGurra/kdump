## Happy polyglot new year (2022)!

Testing new dev platforms is probably good to do from time to time to get some idea of what you are missing out on - and
to see what your strongest biases are.

Seeing an explosion of go based cli and k8s infrastructure components I figured it would be a good excursion to take
from (my case) regular scala/jvm life.

Learning a new platform and/or language works best for me if I my goal is building something that I need (for others I
guess books, videos, discussion groups, reading language spec or something else might work better). In my case I have a
small cli tool called "kdump" which (simplified) dumps all accessible k8s resources from the current kube context into
yaml files on my local file system. Up until this point it has just been a very quick nodejs hack (please dont ask me
why node :D). Basically letting me do diffs, greps, backups, etc on etcd resources before and after certain operations
like helm installs, app rollouts, config updates etc (lets me indirectly figure out how charts, controllers, operators
and other k8s tooling works).

So i rewrote kdump in go. I also rewrote it in Rust (and scala3), because why not (+I have a friend who keeps telling me
how great rust is :) ). The output can be seen here https://github.com/gigurra/kdump (then just click if you want to see
the go or rust implementation).

My thoughts after the experiment?

Overall, in most situations/unless you have very specific control flow and/or complexity requirements, the tooling and
platform simplicity are more important than the programming langauge itself. This was to me quite frustrating, as the
language I am used to (scala) itself is more capable and has simpler solutions for how to model the desired business
logic, but the tooling around it is just.. a nightmare and perfect case of the jvm's "build once * run nowhere".

## [Go](https://github.com/gigurra/kdump/tree/use-go)

While somewhat limiting in terms of raw language capability, its simplicity and incredibly powerful ecosystem for
building, releasing and sharing your work easily makes up for it for desktop/cli apps. Probably would be the same for
consumer device apps, but maybe not so for larger proprietary backend systems. result?: I probably have never been able
to construct an actual useful tool and share it cross platform so easily before. I literally made my application
accessible and deployable to any system without even knowing it (in go * if you have your code in git, preferably semver
tagged, you can deploy and run it anywhere with one command :)).

* language productivity: 6/10. Score based partially on go 1.18 beta incl generics, so I am cheating a little ;) (
  official stable version at time of writing is 1.17). It's a simple language. It's not meant to be super powerful. The
  power comes elsewhere. I did use things like https://github.com/thoas/go-funk though, to at least bring go into the
  modern world a little bit ;P. Had one issue that really bothered me, and that was error handling * but apparently that
  should be improved quite a lot in go 2.0, inspired a bit by Rust ;)
* tooling productivity: 10/10. (the ease of installation, sharing, getting started and deployment is outstanding). One
  installer for either linux and windows. Then go <cmd> <something>.. It was all set up!

* ease to learn language: 10/10
* ease to learn platform: 10/10
* IDE support: 10/10 (i mean, the language is so simple so I didnt have a single issue. IDE was always 1:1 with compiler
  result and other hints were always correct)
* overall enjoyment: 8/10 (I suspect this could increase up as go evolves into a slightly more powerful language)

## [Rust](https://github.com/gigurra/kdump/tree/use-rust)

Totally different beast. Think C++ but make move the default instead of copying. "what do I need from language and
compiler to make *move* the default?"). To my surprise the tooling for building and sharing your work was very simple to
get up and running, similar to go, and made the overall experience really enjoyable. Also dont forget: Rust solves one
problem no other language in this comparison does, and that is being able to code rather large systems without having a
GC :) - meaning there are some realtime situations where the other languages here simply cannot go.

* language productivity: 7/10. Not as limiting as go, but it requires a LOT more brain power from you. Thankfully, even
  though it requires you to be more explicit to pass the compiler's gaze, it doesn't really (to my surprise) result in
  hiding the overall intent that much :) (often a case imo where explicit details can cause you to lose sight of the
  overall goals/purpose of a system,method,class,etc)
* ease to learn language: 6/10 (this is probably Rust's biggest challenge for new people. For old C++ veterans with a
  good grasp on move semantics and who have dabbled into functional coding, it shouldn't be an issue at all. But if you
  decide to make rust the 1st language you learn and expect to be productive within a day or two?... I wouldnt want to
  be in your situation :) )
* ease to learn platform: 8/10 (almost perfect. Reason for not getting a 10 is that the rust installer relies on some
  system dev tools already being installed on your machine, **looking at you windows & ms vc dev tools**, unlike for
  example the go installer. On linux the experience was better/same simplicity as Go)
* IDE support: 7/10 (it's not bad. But it does have some shortcomings when dealing with generic code that involves
  traits, type and lifetime parameters. Nothing major though. I mostly just had to explicitly make the IDE aware of
  types by explicitly writing them out on my helper variables, not relying on IDE's type inference)
* overall enjoyment: 8/10 (while probably not as productive as Go out of the box, it does have its charm and I dont
  think I have ever felt so sure about my application doing the right thing as in Rust * once the app compiled. And
  the "Clippy" linter is awesome for improving your code style and quality!)

## [Scala/jvm base](https://github.com/gigurra/kdump/tree/use-scala)

NOTE: BIAS warning. I have spent roughly 60-70% of my last 10 years of professional development in jvm land :). But
already after having spent a week or so in go/rust land, it seems like our jvm langauge gods/designers have focused way
too long on the language(s) itself and is missing quite a lot of simplicity when it comes to application deployment and
code sharing. Deploying a jvm application is a real hassle, and sharing your code as an opensource library is quite the
dance.

* language productivity: 9/10 (based on scala3 & my subjectively objective experience :D )
* tooling productivity: 4/10. I mean. it is not the worst (*wink* *wink*, looking at you C++ et al, who still havent
  decided on a commonly accepted cross platform package manager)

* ease to learn language: 7/10. Probably slightly easier to learn than rust * but nowhere near Go :).
* ease to learn platform: 4/10. It's bad. Really bad :). Imaginary challenge: start with two compters. One is your dev
  machine with fresh OS installed and one is your target deployment system where other things are already installed (
  possibly other JREs that must not be disturbed). What do you need to install/what commands do you need to develop your
  application on the dev machine and what do you need to do on your target machine? Can you repeat it ? Can you automate
  it? can you maintain it? :D. Honestly.. if it wasnt for container technologies... jvm applications would be dead by
  now ;p.

* IDE support: 8/10 (Slightly more mature than rust land. But Scala is an amazingly powerful language, which is also a
  shortcoming when talking about tooling support. It is HARD for an IDE to cover all corner cases and always give you
  the right information)

* overall enjoyment: Two grades. 8 and 4 :D. Let me explain: If you are already in a corp/enterprise environment where
  all the tooling issues of the JVM have been solved - then it is quite high! But if you are building something from
  scratch, especially if you want do build something simple and possibly open source, then I would honestly not pick a
  jvm based solution (unless you have some very specific reasons for doing so).

## Final biased conclusions

So... what did I learn? well.. I will probably pick either go or rust for any future cli and desktop/consumer device
apps I make. If it is a solo project I might lean towards rust (since it is a little more fun to code in imo), but if I
want to make something shared, then Go is preferred for me. I wish I could say that I would pick scala, but the
ecosystem around it is just too much of a nightmare :(, even though the langauge is the most enjoyable one of them all
to me.

