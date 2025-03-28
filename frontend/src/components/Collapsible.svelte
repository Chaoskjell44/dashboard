<div style="margin-bottom: 8px">
  <div class="inline" class:pointer={!forceAlwaysOpen} on:click={() => toggle(false)}>
    {#if forceAlwaysOpen}
      <i class="fas fa-chevron-right"></i>
    {:else if expanded}
      <i class="{retractIcon}"></i>
    {:else}
      <i class="{expandIcon}"></i>
    {/if}

    <slot name="header"></slot>

    {#if tooltip !== undefined}
      <div style="">
        <Tooltip tip={tooltip} top color="#121212">

          {#if tooltipUrl !== undefined}
            <a href={tooltipUrl} target="_blank">
              <i class="fas fa-circle-info form-label tooltip-icon"></i>
            </a>
          {:else}
            <i class="fas fa-circle-info form-label tooltip-icon"></i>
          {/if}
        </Tooltip>
      </div>
    {/if}

    <hr/>
  </div>

  <div bind:this={content} class="content" class:expanded={expanded}>
    <slot name="content"></slot>
  </div>
</div>

<svelte:window bind:innerWidth />

<script>
    import {afterUpdate, onMount} from "svelte";
    import Tooltip from "svelte-tooltip";
    import { createEventDispatcher } from 'svelte';

    export let retractIcon = "fas fa-minus";
    export let expandIcon = "fas fa-plus";

    export let forceAlwaysOpen = false;
    export let defaultOpen = false;
    export let tooltip = undefined;
    export let tooltipUrl = undefined;

    let expanded = false;
    let showOverflow = true;

    let content;

    let innerWidth;
    $: innerWidth, updateIfExpanded();

    const dispatch = createEventDispatcher();

    export function toggle(force) {
        if (forceAlwaysOpen && !force) {
            return;
        }

        if (expanded) {
            content.style.maxHeight = 0;
        } else {
            updateSize();
        }

        if (expanded) {
            expanded = !expanded;
        } else {
            setTimeout(() => {
                expanded = !expanded;
            }, 300);
        }
    }

    function updateSize() {
        content.style.maxHeight = `100%`;
    }

    function updateIfExpanded() {
        if (expanded) {
            updateSize();
        }
    }

    onMount(() => {
        const observer = new MutationObserver(() => {
            updateIfExpanded();
            setTimeout(updateIfExpanded, 300); // TODO: Move with transition height
        });

        observer.observe(content, { childList: true, subtree: true });

        if (defaultOpen || forceAlwaysOpen) toggle(true);
    });

    afterUpdate(updateIfExpanded);
</script>

<style>
    .content {
        display: flex;
        transition: max-height .3s ease-in-out, margin-top .3s ease-in-out, margin-bottom .3s ease-in-out;
        position: relative;
        overflow: hidden;
        max-height: 0;
    }

    .content.expanded {
        overflow: unset !important;
    }

    .inline {
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 100%;
        gap: 10px;
    }

    hr {
        border-top: 1px solid #777;
        border-bottom: 0;
        border-left: 0;
        border-right: 0;
        width: 100%;
        flex: 1;
    }

    .pointer {
        cursor: pointer;
    }
</style>