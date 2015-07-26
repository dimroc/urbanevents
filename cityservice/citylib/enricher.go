package citylib

type Enricher interface {
	Enrich(g GeoEvent) GeoEvent
}

type broadcastEnricher struct {
	Enrichers []Enricher
}

func NewBroadcastEnricher(enrichers ...Enricher) Enricher {
	return &broadcastEnricher{
		Enrichers: enrichers,
	}
}

func (b *broadcastEnricher) Enrich(g GeoEvent) GeoEvent {
	newGeoevent := g

	for _, enricher := range b.Enrichers {
		newGeoevent = enricher.Enrich(newGeoevent)
	}

	return newGeoevent
}
