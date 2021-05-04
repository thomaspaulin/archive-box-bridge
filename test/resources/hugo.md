---
title: "A Hugo Template"
date: 2021-05-03
description: "Used to test handling of Hugo shortcodes, specifically rel and relref"
tags: ["foo", "bar"]
draft: false
---
# Welcome To My Markdown
I'm some markdown text. I reference [Wikipedia](https://en.wikipedia.org) and I also reference approved [*example*](https://www.example.com) [*sites*](https://www.example.org).

## I Contain Subheaders Too
With associated paragraph text.

I reference [another page using rel]({{< rel "file" >}}) and I reference [one using relref]({{< relref "file" >}}).