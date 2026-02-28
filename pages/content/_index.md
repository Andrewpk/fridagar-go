---
title: "Icelandic Holidays"
description: "A complete listing of public holidays in Iceland — Íslenskir frídagar."
---

{{< holiday-calendar >}}

This site is automatically generated every week at a minimum and will autopublish a new version of the [iCal
compatible ICS file](icelandic_holidays.ics). The calendar will always contain the past year, the current year, and the next year's holidays.

You can use the [ics link here](icelandic_holidays.ics) to add all of the current holidays to your calendar once, or 
you can copy the link and use it to subscribe to the calendar in your preferred calendar app. Subscribing will ensure 
that you always have the most up-to-date holiday information without needing to manually download new files.

<div class="not-prose flex justify-center my-8">
  <div class="relative flex flex-col items-center">
    <button
       id="ics-copy-btn"
       onclick="(function(){
         var url = window.location.origin + window.location.pathname.replace(/\/?$/, '/') + 'icelandic_holidays.ics';
         navigator.clipboard.writeText(url).then(function(){
           var tip = document.getElementById('ics-tooltip');
           tip.classList.remove('opacity-0');
           tip.classList.add('opacity-100');
           setTimeout(function(){ tip.classList.remove('opacity-100'); tip.classList.add('opacity-0'); }, 2000);
         });
       })()"
       class="flex flex-col items-center gap-3 p-8 rounded-2xl border border-neutral-200 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 hover:bg-primary-50 dark:hover:bg-primary-900 hover:border-primary-400 dark:hover:border-primary-500 transition-all duration-200 group cursor-pointer"
       title="Copy calendar URL to clipboard">
      <span class="flex items-center justify-center text-primary-600 dark:text-primary-400 group-hover:scale-110 transition-transform duration-200" style="width:4rem;height:4rem;font-size:3rem;line-height:0;">
        {{< icon "calendar" >}}
      </span>
    </button>
    <span id="ics-tooltip" class="opacity-0 transition-opacity duration-300 mt-3 text-sm font-medium bg-neutral-900 dark:bg-white text-white dark:text-neutral-900 px-4 py-2 rounded-lg pointer-events-none whitespace-nowrap shadow-lg">
      URL copied to clipboard!
    </span>
  </div>
</div>
