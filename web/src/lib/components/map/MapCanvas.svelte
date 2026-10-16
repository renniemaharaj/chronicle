<script lang="ts">
  import { Checkbox, Label } from 'flowbite-svelte';
  import { List, MapPin, Route } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { mapEvents, mapPlaces, recordedVisits } from '../../utils/map-projection';
  import type { Entity } from '../../types/knowledge';
  import MapStage from './MapStage.svelte';
  import MapPlaces from './MapPlaces.svelte';
  import PlaceLocator from './PlaceLocator.svelte';

  let travel = $state(false);
  let throughTime = $state(false);
  let listOpen = $state(true);
  let locating = $state<Entity>();
  const scoped = $derived(ws.startYear !== null || ws.endYear !== null || !!ws.focusedChapterId);
  const events = $derived(mapEvents(ws.project, ws.focusedChapterId, ws.startYear, ws.endYear));
  const visits = $derived(recordedVisits(ws.project, events).filter((event) => !throughTime || (event.time.end?.year ?? event.time.start!.year) <= ws.year));
  const visibleEvents = $derived(travel ? visits : events);
  const places = $derived(ws.project.entities.filter((entity) => entity.kind === 'place' &&
    (!(scoped || travel) || visibleEvents.some((event) => event.locationIds.includes(entity.id)))));
  const points = $derived(mapPlaces(ws.project, visibleEvents, scoped || travel));
</script>

<div class="map-canvas" aria-label="Historical places map">
  <div class="map-toolbar">
    <div class="map-modes" role="group" aria-label="Map display">
      <Button active={!travel} pressed={!travel} onclick={() => { travel = false; }}><MapPin size={14} />Places</Button>
      <Button active={travel} pressed={travel} onclick={() => { travel = true; }}><Route size={14} />Travel</Button>
    </div>
    {#if travel}<Label><Checkbox bind:checked={throughTime} />Through {ws.yearLabel()}</Label>{/if}
    <Button icon active={listOpen} aria-label="Toggle map place list" onclick={() => { listOpen = !listOpen; }}><List size={16} /></Button>
  </div>
  <div class="map-body" class:with-list={listOpen}>
    <div class="map-stage-wrap">
      <MapStage places={points} visits={travel ? visits : []} />
      <div class="map-legend">Modern place references{travel ? ' · Dashed lines: recorded sequence, route unknown' : ''}</div>
      {#if !points.length}<p class="map-empty">{places.length ? 'Choose Locate beside a place to add its map reference.' : 'No places are recorded in this period.'}</p>{/if}
    </div>
    {#if listOpen}<MapPlaces {places} events={visibleEvents} {travel} onlocate={(entity) => { locating = entity; }} />{/if}
    {#if locating}
      {#key locating.id}<PlaceLocator entity={locating} onclose={() => { locating = undefined; }} />{/key}
    {/if}
  </div>
</div>
